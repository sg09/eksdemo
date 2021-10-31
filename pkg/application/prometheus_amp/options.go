package prometheus_amp

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource/amp"
	"fmt"
)

type PrometheusAmpOptions struct {
	application.ApplicationOptions

	AmpEndpoint string
	PushGateway bool
}

func NewOptions() (options *PrometheusAmpOptions, flags cmd.Flags) {
	options = &PrometheusAmpOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "prometheus",
			ServiceAccount: "prometheus",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "v2.30.3",
				Previous: "v2.29.2",
			},
		},
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "push-gateway",
				Description: "enable the Pushgateway",
			},
			Option: &options.PushGateway,
		},
	}
	return
}

func (o *PrometheusAmpOptions) PreInstall() error {
	if o.DryRun {
		o.AmpEndpoint = "<amp-endpoint-goes-here>"
		return nil
	}
	ampGetter := amp.Getter{}

	workspace, err := ampGetter.GetAmpByAlias(fmt.Sprintf("%s-%s", o.ClusterName, ampName))
	if err != nil {
		return fmt.Errorf("failed to lookup AMP endpoint to use in Helm chart for remoteWrite url: %w", err)
	}

	o.AmpEndpoint = *workspace.PrometheusEndpoint

	return nil
}
