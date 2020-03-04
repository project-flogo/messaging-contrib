package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	kafkaConn "github.com/project-flogo/messaging-contrib/kafka/connection"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

var activityMd = activity.ToMetadata(&Input{}, &Output{})

// Activity is a kafka activity
type Activity struct {
	conn  kafkaConn.KafkaConnection
	topic string
}

// New create a new kafka activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	conn, err := coerce.ToConnection(s.Connection)
	if err != nil {
		//ctx.Logger().Errorf("Kafka parameters initialization got error: [%s]", err.Error())
		return nil, err
	}

	act := &Activity{conn: conn.GetConnection().(kafkaConn.KafkaConnection), topic: s.Topic}
	return act, nil
}

// Metadata returns the metadata for the kafka activity
func (*Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements the evaluation of the kafka activity
func (act *Activity) Eval(ctx activity.Context) (done bool, err error) {
	input := &Input{}

	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	if input.Message == "" {
		return false, fmt.Errorf("no message to publish")
	}

	ctx.Logger().Debugf("sending Kafka message")

	msg := &sarama.ProducerMessage{
		Topic: act.topic,
		Value: sarama.StringEncoder(input.Message),
	}

	partition, offset, err := act.conn.Producer().(sarama.SyncProducer).SendMessage(msg)
	if err != nil {
		return false, fmt.Errorf("failed to send Kakfa message for reason [%s]", err.Error())
	}

	output := &Output{}
	output.Partition = partition
	output.OffSet = offset

	if ctx.Logger().DebugEnabled() {
		ctx.Logger().Debugf("Kafka message [%v] sent successfully on partition [%d] and offset [%d]",
			input.Message, partition, offset)
	}

	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}

	return true, nil
}
