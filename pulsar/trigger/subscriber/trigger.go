package subscriber

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
)

const (
	ProcessingModeAsync = "Async"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{})

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

type Trigger struct {
	client   pulsar.Client
	handlers []*Handler
	logger   log.Logger
}
type Handler struct {
	handler                      trigger.Handler
	consumer                     pulsar.Consumer
	done                         chan bool
	asyncMode                    bool
	maxMsgCount, currentMsgCount int
	wg                           sync.WaitGroup
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

	return &Trigger{client: pulsarConn.GetConnection().(pulsar.Client)}, nil
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
	// Init handlers
	for _, handler := range ctx.GetHandlers() {

		s := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), s, true)
		if err != nil {
			return err
		}
		consumeroptions := pulsar.ConsumerOptions{
			Topic:            s.Topic,
			SubscriptionName: s.Subscription,
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

		consumeroptions.MessageChannel = make(chan pulsar.ConsumerMessage)
		consumer, err := t.client.Subscribe(consumeroptions)
		if err != nil {
			return err
		}
		tHandler := &Handler{handler: handler, consumer: consumer, done: make(chan bool)}
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
	for _, handler := range t.handlers {
		go handler.consume()
	}
	t.logger.Info("Trigger Started")
	return nil
}

// Stop implements util.Managed.Stop
func (t *Trigger) Stop() error {
	t.logger.Info("Stopping Trigger")
	for _, handler := range t.handlers {
		// Stop polling
		handler.done <- true
		handler.consumer.Close()
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

func (handler *Handler) consume() {
	handler.handler.Logger().Info("Pulsar Message consumer is started")
	defer handler.handler.Logger().Info("Pulsar Message consumer is stopped")
	for {
		select {
		case msg, ok := <-handler.consumer.Chan():
			if !ok {
				handler.handler.Logger().Error("Error while receiving message")
				time.Sleep(1 * time.Second)
				continue
			}
			// Handle messages concurrently on separate goroutine
			// go handler.handleMessage(msg)
			if handler.asyncMode {
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
		if handler.asyncMode {
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
	out.Properties = msg.Properties()
	out.Topic = msg.Topic()
	handler.handler.Logger().Debugf("Message received [%v]", out.Payload)
	// Do something with the message
	_, err := handler.handler.Handle(context.Background(), out)
	if err == nil {
		// Message processed successfully
		handler.consumer.Ack(msg)
	} else {
		// Failed to process messages
		handler.consumer.Nack(msg)
	}
}
