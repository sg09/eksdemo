package subnet

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "subnet",
			Description: "VPC Subnet",
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{
			ClusterFlagOptional: true,
		},
	}

	return res
}
