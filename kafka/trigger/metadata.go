package kafka

import (
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/connection"
)

type Settings struct {
	Connection connection.Manager `md:"connection,required"`
}
type HandlerSettings struct {
	Topic      string `md:"topic,required"` // The Kafka topic on which to listen for messageS
	Partitions string `md:"partitions"`     // The specific partitions to consume messages from
	Offset     int64  `md:"offset"`         // The offset to use when starting to consume messages, default is set to Newest
}

type Output struct {
	Message string `md:"message"` // The message that was consumed
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message": o.Message,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Message, err = coerce.ToString(values["message"])
	if err != nil {
		return err
	}

	return nil
}
