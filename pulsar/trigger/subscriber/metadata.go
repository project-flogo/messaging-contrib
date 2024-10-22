package subscriber

import (
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/connection"
)

type Settings struct {
	Connection connection.Manager `md:"connection,required"`
}

type HandlerSettings struct {
	Topic                          string `md:"topic"`
	Subscription                   string `md:"subscriptionName,required"`
	SubscriptionType               string `md:"subscriptionType"`
	ProcessingMode                 string `md:"processingMode"`
	InitialPosition                string `md:"initialPosition"`
	Seek                           string `md:"seek"`
	SeekTime                       string `md:"seekTime"`
	LedgerId                       int    `md:"ledgerId"`
	EntryId                        int    `md:"entryId"`
	DLQMaxDeliveries               int    `md:"dlqMaxDeliveries"`
	DLQTopic                       string `md:"dlqTopic"`
	NackRedeliveryDelay            int    `md:"nackRedeliveryDelay"`
	TopicsPattern                  string `md:"topicsPattern"`
	EnableBatchIndexAcknowledgment bool   `md:"enableBatchIndexAcknowledgment"`
	MaxPendingChunkedMessage       int    `md:"maxPendingChunkedMessage"`
	ExpireTimeOfIncompleteChunk    int    `md:"expireTimeOfIncompleteChunk"`
	AutoAckIncompleteChunk         bool   `md:"autoAckIncompleteChunk"`
	ReplicateSubscriptionState     bool   `md:"replicateSubscriptionState"`
}

type Output struct {
	Properties      map[string]string `md:"properties"`
	Payload         interface{}       `md:"payload"`
	Topic           string            `md:"topic"`
	Msgid           string            `md:"msgid"`
	RedeliveryCount int               `md:"redeliveryCount"`
	EntryID         int               `md:"entryid"`
	LedgerID        int               `md:"ledgerid"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.Payload, err = coerce.ToObject(values["payload"])
	if err != nil {
		return err
	}
	o.Properties, err = coerce.ToParams(values["properties"])
	if err != nil {
		return err
	}
	o.Topic, err = coerce.ToString(values["topic"])
	if err != nil {
		return err
	}
	o.Msgid, err = coerce.ToString(values["msgid"])
	if err != nil {
		return err
	}
	o.EntryID, err = coerce.ToInt(values["entryid"])
	if err != nil {
		return err
	}
	o.LedgerID, err = coerce.ToInt(values["ledgerid"])
	if err != nil {
		return err
	}
	o.RedeliveryCount, err = coerce.ToInt(values["redeliveryCount"])
	if err != nil {
		return err
	}
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"payload":         o.Payload,
		"properties":      o.Properties,
		"topic":           o.Topic,
		"msgid":           o.Msgid,
		"redeliveryCount": o.RedeliveryCount,
		"entryid":         o.EntryID,
		"ledgerid":        o.LedgerID,
	}
}
