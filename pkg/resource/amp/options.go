package amp

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

type AmpOptions struct {
	resource.CommonOptions

	Alias string
}

func NewOptions() (options *AmpOptions, flags cmd.Flags) {
	options = &AmpOptions{
		CommonOptions: resource.CommonOptions{
			Name:                "amazon-managed-prometheus",
			ArgumentOptional:    true,
			ClusterFlagDisabled: true,
			DeleteByIdFlag:      true,
		},
	}

	flags = cmd.Flags{}

	return
}

func (o *AmpOptions) SetName(name string) {
	o.Alias = name
}
