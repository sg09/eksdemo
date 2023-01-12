package route_table

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

type RouteTableOptions struct {
	resource.CommonOptions

	VpcId string
}

func newOptions() (options *RouteTableOptions, flags cmd.Flags) {
	options = &RouteTableOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagOptional: true,
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "vpc-id",
				Description: "filter by VPC Id",
				Shorthand:   "V",
			},
			Option: &options.VpcId,
		},
	}
	return
}
