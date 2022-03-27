package security_group

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

type SecurityGroupOptions struct {
	resource.CommonOptions
	NetworkInterfaceId string
	IpAddress          string
}

func NewOptions() (options *SecurityGroupOptions, flags cmd.Flags) {
	options = &SecurityGroupOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagOptional: true,
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "eni-id",
				Description: "filter by Elastic Network Interface Id",
				Shorthand:   "E",
			},
			Option: &options.NetworkInterfaceId,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ip-address",
				Description: "filter by IPv4 Address",
				Shorthand:   "A",
			},
			Option: &options.IpAddress,
		},
	}
	return
}
