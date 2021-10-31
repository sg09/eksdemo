package amp

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

type AmpOptions struct {
	resource.CommonOptions

	AmpName string
}

func NewOptions() (options *AmpOptions, flags cmd.Flags) {
	options = &AmpOptions{
		CommonOptions: resource.CommonOptions{
			Name: "amazon-managed-prometheus",
		},
	}

	flags = cmd.Flags{}

	return
}

func (o *AmpOptions) SetName(name string) {
	o.AmpName = name
}
