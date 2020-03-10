package activity

/*
import (
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/test"
	_ "github.com/project-flogo/messaging-contrib/pulsar/connection"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}
func TestSettings(t *testing.T) {

	settings := &Settings{
		Connection: map[string]interface{}{
			"ref": "github.com/project-flogo/messaging-contrib/pulsar/connection",
			"settings": map[string]interface{}{
				"url": "pulsar://localhost:6605",
			},
		},
		Topic: "sample",
	}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)
	assert.NotNil(t, act)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("payload", "sample")

	_, err = act.Eval(tc)
	assert.Nil(t, err)
}
*/
