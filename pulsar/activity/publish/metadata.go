package publish

import (
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/connection"
)

type Settings struct {
	Connection              connection.Manager `md:"connection"`
	Topic                   string             `md:"topic,required"`
	SendTimeout             int                `md:"sendTimeout"`
	SendMode                string             `md:"sendMode"`
	CompressionType         string             `md:"compressionType"`
	Chunking                bool               `md:"chunking"`
	Batching                bool               `md:"batching"`
	ChunkMaxMessageSize     int                `md:"chunkMaxMessageSize"`
	BatchingMaxMessages     int                `md:"batchingMaxMessages"`
	BatchingMaxSize         int                `md:"batchingMaxSize"`
	BatchingMaxPublishDelay int                `md:"batchingMaxPublishDelay"`
	EnableReplication       bool               `md:"enableReplication"`
	Clusters                string             `md:"clusters"`
}

type Input struct {
	Key        interface{}       `md:"key"`
	Properties map[string]string `md:"properties"`
	Payload    interface{}       `md:"payload"`
}

func (r *Input) FromMap(values map[string]interface{}) (err error) {
	r.Key, err = coerce.ToString(values["key"])
	if err != nil {
		return
	}
	r.Properties, err = coerce.ToParams(values["properties"])
	if err != nil {
		return
	}
	r.Payload, err = coerce.ToAny(values["payload"])
	if err != nil {
		return
	}
	return err
}

func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"payload":    r.Payload,
		"key":        r.Key,
		"properties": r.Properties,
	}
}

// Output of the publish activity
type Output struct {
	Msgid string `md:"msgid"`
}

// FromMap frommap
func (o *Output) FromMap(values map[string]interface{}) (err error) {
	o.Msgid, err = coerce.ToString(values["msgid"])
	if err != nil {
		return
	}
	return
}

// ToMap tomap
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"msgid": o.Msgid,
	}
}
