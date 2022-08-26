package listener

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

type ListenerOptions struct {
	resource.CommonOptions
	LoadBalancerName string
}

func newOptions() (options *ListenerOptions, flags cmd.Flags) {
	options = &ListenerOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "load-balancer",
				Description: "Load Balancer name",
				Shorthand:   "L",
				Required:    true,
			},
			Option: &options.LoadBalancerName,
		},
	}

	return
}
