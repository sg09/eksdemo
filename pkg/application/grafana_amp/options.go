package grafana_amp

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/application/prometheus_amp"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource/amp"
	"fmt"
)

type GrafanaAmpOptions struct {
	application.ApplicationOptions

	AmpEndpoint          string
	DisableIngress       bool
	GrafanaAdminPassword string
	TLSHost              string
}

func NewOptions() (options *GrafanaAmpOptions, flags cmd.Flags) {
	options = &GrafanaAmpOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "grafana",
			ServiceAccount: "grafana",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "8.2.3",
				Previous: "8.1.7",
			},
		},
	}

	flags = cmd.Flags{
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
	return
}

func (o *GrafanaAmpOptions) PreDependencies() error {
	ampGetter := amp.Getter{}

	workspace, err := ampGetter.GetAmpByAlias(fmt.Sprintf("%s-%s", o.ClusterName, prometheus_amp.AmpName))
	if err != nil {
		return fmt.Errorf("failed to lookup AMP endpoint, install prometheus-amp first: %w", err)
	}

	o.AmpEndpoint = *workspace.PrometheusEndpoint

	return nil
}
