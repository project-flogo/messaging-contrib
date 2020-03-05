package activity

import (
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/connection"
)

type Settings struct {
	Connection connection.Manager `md:"connection"`
	Topic      string             `md:"topic,required"`
}

type Input struct {
	Payload string `md:"payload,required"`
}

func (r *Input) FromMap(values map[string]interface{}) error {
	var err error
	r.Payload, err = coerce.ToString(values["payload"])

	return err
}

func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"payload": r.Payload,
	}
}
