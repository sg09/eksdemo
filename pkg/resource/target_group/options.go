package target_group

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
	"strings"

	"github.com/spf13/cobra"
)

type TargeGroupOptions struct {
	resource.CommonOptions

	LoadBalancerName string
	Protocol         string
	TargetType       string
}

func newOptions() (options *TargeGroupOptions, createFlags, getFlags cmd.Flags) {
	options = &TargeGroupOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagOptional: true,
		},
		Protocol:   "http",
		TargetType: "instance",
	}

	createFlags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "protocol",
				Description: "protocol for routing traffic to the targets",
				Shorthand:   "p",
				Validate: func(cmd *cobra.Command, args []string) error {
					options.Protocol = strings.ToUpper(options.Protocol)
					return nil
				},
			},
			Choices: []string{"http", "https", "tcp", "tls", "udp", "tcp_udp", "geneve"},
			Option:  &options.Protocol,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "target-type",
				Description: "target type",
				Shorthand:   "t",
			},
			Choices: []string{"ip", "instance", "lambda", "alb"},
			Option:  &options.TargetType,
		},
	}

	getFlags = cmd.Flags{
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
