package kube_prometheus

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
)

type KubePrometheusOptions struct {
	*application.ApplicationOptions
	GrafanaAdminPassword string
	IngressHost          string
}

func addOptions(a *application.Application) *application.Application {
	options := &KubePrometheusOptions{
		ApplicationOptions: &application.ApplicationOptions{
			DisableServiceAccountFlag: true,
			Namespace:                 "monitoring",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "34.10.0",
				Latest:        "v0.55.0",
				PreviousChart: "34.10.0",
				Previous:      "v0.55.0",
			},
		},
	}
	a.Options = options

	a.Flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "grafana-pass",
				Description: "Grafana admin password (required)",
				Required:    true,
			},
			Option: &options.GrafanaAdminPassword,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ingress-host",
				Description: "hostname for Ingress with TLS (requires ACM cert, AWS LB Controller and ExternalDNS)",
				Shorthand:   "I",
			},
			Option: &options.IngressHost,
		},
	}
	return a
}
