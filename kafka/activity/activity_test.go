package kafka

import (
	"encoding/json"
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/connection"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}

const settingsConfig string = `{
	"connection": {
		"ref":"github.com/project-flogo/messaging-contrib/kafka/connection",
		"settings":{
			"brokerUrls":"localhost:9092"
		}
	}
	}`

func TestPlain(t *testing.T) {

	var s map[string]*connection.Config
	err := json.Unmarshal([]byte(settingsConfig), &s)
	assert.Nil(t, err)
	settings := &Settings{}
	settings.Connection, err = connection.NewManager(s["connection"])
	assert.Nil(t, err)
	settings.Topic = "sample"

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("message", "Hello")

	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)
}
