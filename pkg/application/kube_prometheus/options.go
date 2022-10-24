package kube_prometheus

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
)

type KubePrometheusOptions struct {
	*application.ApplicationOptions
	GrafanaAdminPassword string
}

func newOptions() (options *KubePrometheusOptions, flags cmd.Flags) {
	options = &KubePrometheusOptions{
		ApplicationOptions: &application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "41.5.1",
				Latest:        "v0.60.1",
				PreviousChart: "34.10.0",
				Previous:      "v0.55.0",
			},
			DisableServiceAccountFlag:    true,
			ExposeIngressAndLoadBalancer: true,
			Namespace:                    "monitoring",
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "grafana-pass",
				Description: "grafana admin password",
				Required:    true,
				Shorthand:   "P",
			},
			Option: &options.GrafanaAdminPassword,
		},
	}
	return
}
