package kubecost_eks_amp

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/amp"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
)

const AmpAliasSuffix = "kubecost-amp"

type KubecostEksAmpOptions struct {
	application.ApplicationOptions

	AmpEndpoint string
	AmpId       string
	*amp.AmpOptions
}

func newOptions() (options *KubecostEksAmpOptions, flags cmd.Flags) {
	options = &KubecostEksAmpOptions{
		ApplicationOptions: application.ApplicationOptions{
			ExposeIngressAndLoadBalancer: true,
			Namespace:                    "kubecost",
			ServiceAccount:               "kubecost-cost-analyzer",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.97.0",
				Latest:        "1.97.0",
				PreviousChart: "1.97.0",
				Previous:      "1.97.0",
			},
		},
		AmpOptions: &amp.AmpOptions{
			CommonOptions: resource.CommonOptions{
				Name: "kubecost-amazon-managed-prometheus",
			},
		},
	}

	// flags = cmd.Flags{
	// 	&cmd.StringFlag{
	// 		CommandFlag: cmd.CommandFlag{
	// 			Name:        "provider-version",
	// 			Description: "version of provider-aws",
	// 		},
	// 		Option: &options.ProviderVersion,
	// 	},
	// }
	return
}

func (o *KubecostEksAmpOptions) PreDependencies(application.Action) error {
	o.AmpOptions.Alias = fmt.Sprintf("%s-%s", o.ClusterName, AmpAliasSuffix)
	return nil
}

func (o *KubecostEksAmpOptions) PreInstall() error {
	o.AmpEndpoint = "<-amp_endpoint_url_will_go_here->"
	o.AmpId = "<-amp_id_will_go_here->"

	workspace, err := amp.NewGetter(aws.NewAMPClient()).GetAmpByAlias(fmt.Sprintf("%s-%s", o.ClusterName, AmpAliasSuffix))
	if err != nil {
		if o.DryRun {
			return nil
		}
		return fmt.Errorf("failed to lookup AMP to use in Helm chart values file: %w", err)
	}

	o.AmpEndpoint = awssdk.ToString(workspace.PrometheusEndpoint)
	o.AmpId = awssdk.ToString(workspace.WorkspaceId)

	return nil
}
