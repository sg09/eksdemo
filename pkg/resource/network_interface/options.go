package network_interface

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

type NetworkInterfaceOptions struct {
	resource.CommonOptions
	InstanceId      string
	IpAddress       string
	SecurityGroupId string
}

func NewOptions() (options *NetworkInterfaceOptions, flags cmd.Flags) {
	options = &NetworkInterfaceOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagOptional: true,
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "instance-id",
				Description: "filter by Instance Id",
				Shorthand:   "I",
			},
			Option: &options.InstanceId,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ip-address",
				Description: "filter by IPv4 Address",
				Shorthand:   "A",
			},
			Option: &options.IpAddress,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "security-group-id",
				Description: "filter by Security Group Id",
				Shorthand:   "S",
			},
			Option: &options.SecurityGroupId,
		},
	}
	return
}
