package ec2_instance

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "ec2-instance",
			Description: "EC2 Instance",
			Aliases:     []string{"ec2-instances", "ec2", "instances", "instance"},
			Args:        []string{"ID"},
		},

		Getter: &Getter{},

		Manager: &Manager{},

		Options: &resource.CommonOptions{
			ClusterFlagOptional: true,
		},
	}

	return res
}
