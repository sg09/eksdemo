package kube_prometheus

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
)

type KubePrometheusOptions struct {
	*application.ApplicationOptions
	GrafanaAdminPassword string
	DisableIngress       bool
	TLSHost              string
}

func addOptions(a *application.Application) *application.Application {
	options := &KubePrometheusOptions{
		ApplicationOptions: &application.ApplicationOptions{
			DisableServiceAccountFlag: true,
			Namespace:                 "monitoring",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "v0.52.0",
				Previous: "v0.51.2",
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
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "disable-ingress",
				Description: "don't create Ingress for Grafana",
			},
			Option: &options.DisableIngress,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "tls-host",
				Description: "FQDN of host to secure with TLS (requires ExternalDNS for cert discovery) ",
			},
			Option: &options.TLSHost,
		},
	}
	return a
}
