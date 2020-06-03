package function

import "github.com/project-flogo/core/data/coerce"

type Output struct {
	Message []byte `md:"message"`
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Message, err = coerce.ToBytes(values["message"])
	if err != nil {
		return err
	}

	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message": o.Message,
	}
}

type Reply struct {
	Out interface{} `md:"out"`
}

func (r *Reply) FromMap(values map[string]interface{}) error {

	r.Out = values["out"]

	return nil
}

func (r *Reply) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"out": r.Out,
	}
}
