package log_group

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "log-group",
			Description: "CloudWatch Log Group",
			Aliases:     []string{"log-groups", "loggroup", "lg"},
			Args:        []string{"NAME_PREFIX"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	return res
}
