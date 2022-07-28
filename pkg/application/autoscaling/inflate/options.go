package inflate

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
)

type InflateOptions struct {
	application.ApplicationOptions

	Replicas int
}

func NewOptions() (options *InflateOptions, flags cmd.Flags) {
	options = &InflateOptions{
		ApplicationOptions: application.ApplicationOptions{
			DisableServiceAccountFlag: true,
			DisableVersionFlag:        true,
			Namespace:                 "inflate",
		},
		Replicas: 0,
	}

	flags = cmd.Flags{
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "replicas",
				Description: "number of replicas for the deployment (default 0)",
			},
			Option: &options.Replicas,
		},
	}
	return
}
