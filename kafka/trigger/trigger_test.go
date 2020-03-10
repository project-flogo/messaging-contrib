package kafka

import (
	"encoding/json"
	"testing"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/support/test"
	"github.com/project-flogo/core/trigger"
	_ "github.com/project-flogo/messaging-contrib/kafka/connection"
	"github.com/stretchr/testify/assert"
)

const testConfig string = `{
	"id": "flogo-kafka",
	"ref": "github.com/project-flogo/contrib/trigger/kafka",
	"settings":{
		"connection": {
			"ref": "github.com/project-flogo/messaging-contrib/kafka/connection",
			"settings":{
			  "brokerUrls": "localhost:9092"
			}
		}
	},
	"handlers": [
	  {
			"action":{
				"id":"dummy"
			},
			"settings": {
		  		"topic": "sample"
			}
	  }
	]

  }`

func TestKafkaTrigger_Initialize(t *testing.T) {
	f := &Factory{}

	config := &trigger.Config{}
	err := json.Unmarshal([]byte(testConfig), config)
	assert.Nil(t, err)

	actions := map[string]action.Action{"dummy": test.NewDummyAction(func() {
		//do nothing
	})}

	trg, err := test.InitTrigger(f, config, actions)
	assert.Nil(t, err)
	assert.NotNil(t, trg)

	err = trg.Start()
	assert.Nil(t, err)
	err = trg.Stop()
	assert.Nil(t, err)

}
