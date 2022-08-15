package target_group

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

type TargeGroupOptions struct {
	resource.CommonOptions
	LoadBalancerName string
}

func newOptions() (options *TargeGroupOptions, flags cmd.Flags) {
	options = &TargeGroupOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagOptional: true,
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "load-balancer",
				Description: "filter by Load Balancer name",
				Shorthand:   "L",
			},
			Option: &options.LoadBalancerName,
		},
	}

	return
}
