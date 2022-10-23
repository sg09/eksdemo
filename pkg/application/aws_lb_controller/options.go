package aws_lb_controller

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
)

type AWSLBControllerOptions struct {
	application.ApplicationOptions

	Default bool
}

func newOptions() (options *AWSLBControllerOptions, flags cmd.Flags) {
	options = &AWSLBControllerOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "awslb",
			ServiceAccount: "aws-load-balancer-controller",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.4.5",
				Latest:        "v2.4.4",
				PreviousChart: "1.4.4",
				Previous:      "v2.4.3",
			},
		},
		Default: false,
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "default",
				Description: "set as the default IngressClass for the cluster",
			},
			Option: &options.Default,
		},
	}

	return
}
