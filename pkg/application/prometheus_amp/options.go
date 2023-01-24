package prometheus_amp

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource/amp_workspace"
	"fmt"
)

const AmpAliasSuffix = "prometheus-amp"

type PrometheusAmpOptions struct {
	application.ApplicationOptions

	AmpEndpoint string
	PushGateway bool
	*amp_workspace.AmpWorkspaceOptions
}

func NewOptions() (options *PrometheusAmpOptions, flags cmd.Flags) {
	options = &PrometheusAmpOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "41.5.1",
				Latest:        "v0.60.1",
				PreviousChart: "34.10.0",
				Previous:      "v0.55.0",
			},
			Namespace:      "prometheus",
			ServiceAccount: "prometheus",
		},
	}
	return
}

func (o *PrometheusAmpOptions) PreDependencies(application.Action) error {
	o.AmpWorkspaceOptions.Alias = fmt.Sprintf("%s-%s", o.ClusterName, AmpAliasSuffix)
	return nil
}

func (o *PrometheusAmpOptions) PreInstall() error {
	if o.DryRun {
		o.AmpEndpoint = "<amp-endpoint-goes-here>"
		return nil
	}
	ampGetter := amp_workspace.NewGetter(aws.NewAMPClient())

	workspace, err := ampGetter.GetAmpByAlias(fmt.Sprintf("%s-%s", o.ClusterName, AmpAliasSuffix))
	if err != nil {
		return fmt.Errorf("failed to lookup AMP endpoint to use in Helm chart for remoteWrite url: %w", err)
	}

	o.AmpEndpoint = *workspace.Workspace.PrometheusEndpoint

	return nil
}
