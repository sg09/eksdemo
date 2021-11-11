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
			ClusterFlagDisabled: true,
			DeleteById:          true,
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "alias",
				Description: "workspace alias",
			},
			Option: &options.Alias,
		},
	}

	return
}

func (o *AmpOptions) SetName(name string) {
	o.Alias = name
}
