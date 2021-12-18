package dns_record

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

type DnsRecordOptions struct {
	resource.CommonOptions

	ZoneName string
}

func NewOptions() (options *DnsRecordOptions, flags cmd.Flags) {
	options = &DnsRecordOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "zone",
				Description: "hosted zone name",
				Shorthand:   "z",
				Required:    true,
			},
			Option: &options.ZoneName,
		},
	}
	return
}
