package nginx

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
)

type NginxOptions struct {
	application.ApplicationOptions

	Replicas int
}

func newOptions() (options *NginxOptions, flags cmd.Flags) {
	options = &NginxOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "ingress-nginx",
			ServiceAccount: "ingress-nginx",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "4.2.0",
				Latest:        "v1.3.0",
				PreviousChart: "4.1.4",
				Previous:      "v1.2.1",
			},
		},
		Replicas: 1,
	}

	flags = cmd.Flags{
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "replicas",
				Description: "number of replicas for the deployment",
			},
			Option: &options.Replicas,
		},
	}
	return
}
