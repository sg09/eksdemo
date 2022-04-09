package prometheus_amp

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource/amp"
	"fmt"
)

const AmpAliasSuffix = "prometheus-amp"

type PrometheusAmpOptions struct {
	application.ApplicationOptions

	AmpEndpoint string
	PushGateway bool
	*amp.AmpOptions
}

func NewOptions() (options *PrometheusAmpOptions, flags cmd.Flags) {
	options = &PrometheusAmpOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "prometheus",
			ServiceAccount: "prometheus",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "v0.55.1",
				Previous: "v0.54.1",
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

func (o *PrometheusAmpOptions) PreDependencies(application.Action) error {
	o.AmpOptions.Alias = fmt.Sprintf("%s-%s", o.ClusterName, AmpAliasSuffix)
	return nil
}

func (o *PrometheusAmpOptions) PreInstall() error {
	if o.DryRun {
		o.AmpEndpoint = "<amp-endpoint-goes-here>"
		return nil
	}
	ampGetter := amp.Getter{}

	workspace, err := ampGetter.GetAmpByAlias(fmt.Sprintf("%s-%s", o.ClusterName, AmpAliasSuffix))
	if err != nil {
		return fmt.Errorf("failed to lookup AMP endpoint to use in Helm chart for remoteWrite url: %w", err)
	}

	o.AmpEndpoint = *workspace.PrometheusEndpoint

	return nil
}
