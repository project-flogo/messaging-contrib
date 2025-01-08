package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/project-flogo/core/app"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/engine"
	cnn "github.com/project-flogo/core/support/connection"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/trace"
	"github.com/project-flogo/core/trigger"
	connection "github.com/project-flogo/messaging-contrib/pulsar/connection"
)

const (
	ProcessingModeAsync = "Async"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{})
var EventBasedFlowControl bool

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

type Trigger struct {
	connMgr   connection.PulsarConnManager
	pulsarCnn cnn.Manager
	handlers  []*Handler
	logger    log.Logger
}
type Handler struct {
	handler                      trigger.Handler
	consumer                     pulsar.Consumer
	done                         chan bool
	asyncMode                    bool
	maxMsgCount, currentMsgCount int
	wg                           sync.WaitGroup
	consumerOpts                 pulsar.ConsumerOptions
	seekBy                       string
	seekMessageID                pulsar.MessageID
	seekTime                     time.Time
}

type Factory struct {
}

func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}
	pulsarConn, err := coerce.ToConnection(s.Connection)
	if err != nil {
		return nil, err
	}
	connMgr := pulsarConn.GetConnection().(connection.PulsarConnManager)
	return &Trigger{connMgr: connMgr, pulsarCnn: pulsarConn}, nil
}

func (f *Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// Metadata implements trigger.Trigger.Metadata
func (t *Trigger) Metadata() *trigger.Metadata {
	return triggerMd
}

func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	t.logger = ctx.Logger()
	//Check if flow control is enabled
	EventBasedFlowControl = app.EnableFlowControl()
	// Init handlers
	for _, handler := range ctx.GetHandlers() {

		s := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), s, true)
		if err != nil {
			return err
		}
		var hostName string
		hostName, err = os.Hostname()
		if err != nil {
			hostName = fmt.Sprintf("%d", time.Now().UnixMilli())
		}

		if s.Topic == "" && s.TopicsPattern == "" {
			return fmt.Errorf("Topic is required. Please provide either the Topic or Topics Pattern.")
		}

		consumeroptions := pulsar.ConsumerOptions{
			Topic:                          s.Topic,
			TopicsPattern:                  s.TopicsPattern,
			SubscriptionName:               s.Subscription,
			Name:                           fmt.Sprintf("%s-%s-%s-%s", engine.GetAppName(), engine.GetAppVersion(), handler.Name(), hostName),
			EnableBatchIndexAcknowledgment: s.EnableBatchIndexAcknowledgment,
			MaxPendingChunkedMessage:       s.MaxPendingChunkedMessage,
			ExpireTimeOfIncompleteChunk:    time.Duration(s.ExpireTimeOfIncompleteChunk) * time.Second,
			AutoAckIncompleteChunk:         s.AutoAckIncompleteChunk,
			ReplicateSubscriptionState:     s.ReplicateSubscriptionState,
		}

		if s.NackRedeliveryDelay != 0 {
			consumeroptions.NackRedeliveryDelay = time.Duration(s.NackRedeliveryDelay) * time.Second
		}

		switch s.SubscriptionType {
		case "Exclusive":
			consumeroptions.Type = pulsar.Exclusive
		case "Shared":
			consumeroptions.Type = pulsar.Shared
		case "Failover":
			consumeroptions.Type = pulsar.Failover
		case "KeyShared":
			consumeroptions.Type = pulsar.KeyShared
		default:
			consumeroptions.Type = pulsar.Exclusive
		}
		if s.DLQTopic != "" {
			policy := pulsar.DLQPolicy{
				MaxDeliveries:   uint32(s.DLQMaxDeliveries),
				DeadLetterTopic: s.DLQTopic,
			}
			consumeroptions.DLQ = &policy
		}
		if s.InitialPosition == "Latest" {
			consumeroptions.SubscriptionInitialPosition = pulsar.SubscriptionPositionLatest
		} else {
			consumeroptions.SubscriptionInitialPosition = pulsar.SubscriptionPositionEarliest
		}
		// Seek Functionality
		var seekMessageID pulsar.MessageID
		var seekTime time.Time
		if s.Seek == "Seek By MessageID" {
			seekMessageID = pulsar.NewMessageID(int64(s.LedgerId), int64(s.EntryId), 0, 0)
		} else if s.Seek == "Seek By Timestamp" {
			seekTime, err = time.Parse(time.RFC3339, s.SeekTime)
			if err != nil {
				return fmt.Errorf("Error while parsing seek timestamp : ", err.Error())
			}
		}
		consumeroptions.MessageChannel = make(chan pulsar.ConsumerMessage)
		var consumer pulsar.Consumer

		tHandler := &Handler{handler: handler, consumer: consumer, done: make(chan bool), consumerOpts: consumeroptions, seekBy: s.Seek, seekMessageID: seekMessageID, seekTime: seekTime}
		tHandler.asyncMode = s.ProcessingMode == ProcessingModeAsync
		tHandler.maxMsgCount = getMaxMessageCount()
		tHandler.wg = sync.WaitGroup{}
		t.handlers = append(t.handlers, tHandler)
	}

	return nil
}

func getMaxMessageCount() int {
	if engine.GetRunnerType() == engine.ValueRunnerTypePooled {
		return engine.GetRunnerWorkers()
	}
	// For DIRECT mode
	return 200
}

// Start implements util.Managed.Start
func (t *Trigger) Start() error {
	t.logger.Info("Starting Trigger")
	t.connMgr = t.pulsarCnn.GetConnection().(connection.PulsarConnManager)
	for _, handler := range t.handlers {

		go handler.consume(t.connMgr)
	}
	t.logger.Info("Trigger Started")
	return nil
}

// Stop implements util.Managed.Stop
func (t *Trigger) Stop() error {
	t.logger.Info("Stopping Trigger")
	for _, handler := range t.handlers {
		// Stop polling
		if handler.consumer != nil {
			handler.done <- true
			handler.consumer.Close()
		}
	}
	t.logger.Info("Trigger Stopped")
	return nil
}

func (t *Trigger) Resume() error {
	t.logger.Info("Resuming Trigger")
	err := t.Start()
	if err == nil {
		t.logger.Info("Trigger Resumed")
	}
	return err
}

func (t *Trigger) Pause() error {
	for _, handler := range t.handlers {
		handler.done <- true
	}
	t.logger.Info("Trigger Paused")
	return nil
}

func (handler *Handler) consume(connMgr connection.PulsarConnManager) {

	var err error

	for handler.consumer == nil {
		handler.handler.Logger().Debugf("Attempting subscriber creation for handler %v", handler.handler.Name())
		handler.consumer, err = connMgr.GetSubscriber(handler.consumerOpts)
		if err != nil {
			handler.handler.Logger().Errorf("%v", err)

			handler.handler.Logger().Infof("Retrying connection after 60 seconds")
			time.Sleep(60 * time.Second)
		}
	}
	// Seek
	if handler.seekBy == "Seek By MessageID" {
		err = handler.consumer.Seek(handler.seekMessageID)
		if err != nil {
			handler.handler.Logger().Errorf("%v", err)
		}
	} else if handler.seekBy == "Seek By Timestamp" {
		err = handler.consumer.SeekByTime(handler.seekTime)
		if err != nil {
			handler.handler.Logger().Errorf("%v", err)
		}
	}

	defer handler.handler.Logger().Info("Pulsar Message consumer is stopped")
	handler.handler.Logger().Info("Pulsar Message consumer is started")
	for {
		select {
		case msg, ok := <-handler.consumer.Chan():
			if !ok {
				handler.handler.Logger().Error("Error while receiving message")
				time.Sleep(1 * time.Second)
				continue
			}

			//If flow control is enabled
			// use engine's flow limit control
			if EventBasedFlowControl {
				go handler.handleMessage(msg)
			} else if handler.asyncMode {
				// Handle messages concurrently on separate goroutine
				// go handler.handleMessage(msg)
				handler.wg.Add(1)
				handler.currentMsgCount++
				go handler.handleMessage(msg)
				if handler.currentMsgCount >= handler.maxMsgCount {
					handler.handler.Logger().Infof("Total messages received are equal or more than maximum threshold [%d]. Blocking message handler.", handler.maxMsgCount)
					handler.wg.Wait()
					// reset count
					handler.currentMsgCount = 0
					handler.handler.Logger().Info("All received messages are processed. Unblocking message handler.")
				}
			} else {
				handler.handleMessage(msg)
			}
		case <-handler.done:
			return
		}
	}
}

func (handler *Handler) handleMessage(msg pulsar.ConsumerMessage) {
	defer func() {
		if !EventBasedFlowControl && handler.asyncMode {
			handler.wg.Done()
			handler.currentMsgCount--
		}
	}()
	handler.handler.Logger().Debugf("Message received - %s", msg.ID())
	out := &Output{}
	if handler.handler.Settings()["format"] != nil &&
		handler.handler.Settings()["format"].(string) == "JSON" {
		var obj interface{}
		err := json.Unmarshal(msg.Payload(), &obj)
		if err != nil {
			handler.handler.Logger().Errorf("Pulsar consumer, configured to receive JSON formatted messages, was unable to parse message: [%v]", msg.Payload())
			handler.consumer.Nack(msg)
			return
		}
		out.Payload = obj
	} else {
		out.Payload = string(msg.Payload())
	}

	ctx := context.Background()
	if trace.Enabled() {
		tc, _ := trace.GetTracer().Extract(trace.TextMap, msg.Properties())
		if tc != nil {
			ctx = trace.AppendTracingContext(ctx, tc)
		}
	}
	out.Properties = msg.Properties()
	out.Topic = msg.Topic()
	out.RedeliveryCount = int(msg.RedeliveryCount())
	msgID := msg.ID()
	if msgID != nil {
		out.Msgid = fmt.Sprintf("%x", msgID.Serialize())
	}
	handler.handler.Logger().Debugf("Message received [%v] with msgID [%v]", out.Payload, out.Msgid)
	// Log EntryID,LedgerID and PartitionIdx logs
	handler.handler.Logger().Debugf("LedgerId  [%v] : EntryID [%v] : PartitionIdx [%v]", msgID.LedgerID(), msgID.EntryID(), msgID.PartitionIdx())
	// Do something with the message
	if out.Msgid != "" {
		ctx = trigger.NewContextWithEventId(ctx, out.Msgid)
	}
	out.EntryID = int(msgID.EntryID())
	out.LedgerID = int(msgID.LedgerID())
	attrs, err := handler.handler.Handle(ctx, out)
	if err == nil {
		// Message processed successfully
		if attrs[" _nack"] != nil && attrs[" _nack"] == true {
			handler.consumer.Nack(msg)
		} else {
			handler.consumer.Ack(msg)
		}
	} else {
		// Failed to process messages
		handler.consumer.Nack(msg)
	}
}
