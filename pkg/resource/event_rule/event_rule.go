package event_rule

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "event-rule",
			Description: "EventBridge Rule",
			Aliases:     []string{"event-rules", "event-rule", "eventrules", "eventrule"},
			Args:        []string{"NAME_PREFIX"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	return res
}
