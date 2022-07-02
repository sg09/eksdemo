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
	GrafanaAdminPassword string
	IngressHost          string
}

func NewOptions() (options *GrafanaAmpOptions, flags cmd.Flags) {
	options = &GrafanaAmpOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "grafana",
			ServiceAccount: "grafana",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "34.10.0",
				Latest:        "8.5.0",
				PreviousChart: "34.10.0",
				Previous:      "8.5.0",
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
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ingress-host",
				Description: "hostname for Ingress with TLS (requires ACM cert, AWS LB Controller and ExternalDNS)",
				Shorthand:   "I",
			},
			Option: &options.IngressHost,
		},
	}
	return
}

func (o *GrafanaAmpOptions) PreDependencies(action application.Action) error {
	if action == application.Uninstall {
		return nil
	}

	ampGetter := amp.Getter{}

	workspace, err := ampGetter.GetAmpByAlias(fmt.Sprintf("%s-%s", o.ClusterName, prometheus_amp.AmpAliasSuffix))
	if err != nil {
		return fmt.Errorf("failed to lookup AMP endpoint, install prometheus-amp first: %w", err)
	}

	o.AmpEndpoint = *workspace.PrometheusEndpoint

	return nil
}
