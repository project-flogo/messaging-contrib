package trigger

/*
import (
	"encoding/json"
	"testing"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/support/test"
	"github.com/project-flogo/core/trigger"
	_ "github.com/project-flogo/messaging-contrib/pulsar/connection"
	"github.com/stretchr/testify/assert"
)

const testConfig string = `{
	"id": "flogo-pulsar",
	"ref": " github.com/project-flogo/messaging-contrib/pulsar/trigger",
	"settings": {
	  "connection": {
		  "ref": "github.com/project-flogo/messaging-contrib/pulsar/connection",
		  "settings":{
			"url": "pulsar://localhost:6605"
		  }
	  }
	},
	"handlers": [
	  {
			"action":{
				"id":"dummy"
			},
			"settings": {
			  "topic": "sample",
			  "subscription":"sample-1"
			}
	  }
	]

  }`

func TestPulsarTrigger_Initialize(t *testing.T) {
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

}
*/
