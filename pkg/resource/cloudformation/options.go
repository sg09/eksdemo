package cloudformation

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

type CloudFormationOptions struct {
	resource.CommonOptions
}

func NewOptions() (options *CloudFormationOptions, flags cmd.Flags) {
	options = &CloudFormationOptions{
		CommonOptions: resource.CommonOptions{
			Name:                "cloudformation",
			ClusterFlagDisabled: true,
			ClusterFlagOptional: true,
		},
	}

	flags = cmd.Flags{}

	return
}
