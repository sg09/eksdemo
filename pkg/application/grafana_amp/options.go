package grafana_amp

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/application/prometheus_amp"
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource/amp"
	"fmt"
)

type GrafanaAmpOptions struct {
	application.ApplicationOptions

	AmpEndpoint          string
	GrafanaAdminPassword string
}

func newOptions() (options *GrafanaAmpOptions, flags cmd.Flags) {
	options = &GrafanaAmpOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "39.6.0",
				Latest:        "9.0.5",
				PreviousChart: "34.10.0",
				Previous:      "8.5.0",
			},
			ExposeIngressAndLoadBalancer: true,
			Namespace:                    "grafana",
			ServiceAccount:               "grafana",
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "grafana-pass",
				Description: "Grafana admin password",
				Required:    true,
				Shorthand:   "P",
			},
			Option: &options.GrafanaAdminPassword,
		},
	}
	return
}

func (o *GrafanaAmpOptions) PreDependencies(action application.Action) error {
	if action == application.Uninstall {
		return nil
	}

	ampGetter := amp.NewGetter(aws.NewAMPClient())

	workspace, err := ampGetter.GetAmpByAlias(fmt.Sprintf("%s-%s", o.ClusterName, prometheus_amp.AmpAliasSuffix))
	if err != nil {
		return fmt.Errorf("failed to lookup AMP endpoint, install prometheus-amp first: %w", err)
	}

	o.AmpEndpoint = *workspace.PrometheusEndpoint

	return nil
}
