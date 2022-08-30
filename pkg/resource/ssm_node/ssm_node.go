package ssm_node

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "ssm-node",
			Description: "SSM Managed Node",
			Aliases:     []string{"ssm"},
			Args:        []string{"INSTANCE_ID"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	return res
}
