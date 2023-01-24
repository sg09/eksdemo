package amp_workspace

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

type AmpWorkspaceOptions struct {
	resource.CommonOptions

	Alias string
}

func NewOptions() (options *AmpWorkspaceOptions, flags cmd.Flags) {
	options = &AmpWorkspaceOptions{
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

func (o *AmpWorkspaceOptions) SetName(name string) {
	o.Alias = name
}
