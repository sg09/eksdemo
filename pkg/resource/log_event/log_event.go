package log_event

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "log-event",
			Description: "CloudWatch Log Events",
			Aliases:     []string{"log-events", "logevents", "logs"},
			Args:        []string{"LOG_STREAM_NAME"},
		},

		Getter: &Getter{},
	}

	res.Options, res.GetFlags = newOptions()

	return res
}
