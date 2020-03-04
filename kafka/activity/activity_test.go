package kafka

/*
import (
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}

func TestPlain(t *testing.T) {
	settings := &Settings{
		Connection: map[string]interface{}{
			"ref": "github.com/project-flogo/messaging-contrib/kafka/connection",
			"settings": map[string]interface{}{
				"brokerUrls": "localhost:9092",
			},
		},
		Topic: "sample",
	}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("message", "Hello")

	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)
}
*/
