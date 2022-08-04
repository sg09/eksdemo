package dns_record

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

type DnsRecordOptions struct {
	resource.CommonOptions

	All      bool
	ZoneName string
}

func NewOptions() (options *DnsRecordOptions, deleteFlags cmd.Flags, getFlags cmd.Flags) {
	options = &DnsRecordOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	commonFlags := cmd.Flags{
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

	deleteFlags = append(cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "all",
				Description: "delete all records (including TXT) that start with NAME ",
				Shorthand:   "A",
			},
			Option: &options.All,
		},
	}, commonFlags...)

	getFlags = append(cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "all",
				Description: "show all records (including SOA and TXT records)",
				Shorthand:   "A",
			},
			Option: &options.All,
		},
	}, commonFlags...)

	return
}
