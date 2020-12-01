package acknowledge

import (
	"github.com/project-flogo/core/activity"
)

var activityMd = activity.ToMetadata()

type MyActivity struct {
}

func init() {
	_ = activity.Register(&MyActivity{}, New)
}

func New(ctx activity.InitContext) (activity.Activity, error) {
	return &MyActivity{}, nil
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {
	attrs := make(map[string]interface{})
	attrs["pulsarack"] = true
	context.ActivityHost().Reply(attrs, nil)
	context.Logger().Info("Notifying Pulsar Subscriber to ack the message")
	return true, nil
}
