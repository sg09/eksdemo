package container_insights

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/kustomize"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:     https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Container-Insights-setup-EKS-quickstart.html
// GitHub:   https://github.com/aws-samples/amazon-cloudwatch-container-insights, see Releases for versions
// Manifest: https://github.com/aws-samples/amazon-cloudwatch-container-insights/blob/master/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/quickstart/cwagent-fluent-bit-quickstart.yaml
// Repo:     https://gallery.ecr.aws/cloudwatch-agent/cloudwatch-agent
// Version:  Latest is 1.247348.0b251302 aka k8s/1.3.8 (as of 9/23/21)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "container-insights",
			Description: "CloudWatch Container Insights",
			Aliases:     []string{"ci"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"},
			}),
		},

		Installer: &kustomize.KustomizeInstaller{
			ResourceTemplate: &template.TextTemplate{
				Template: containerInsightsManifestTemplate,
			},
			KustomizeTemplate: &template.TextTemplate{
				Template: kustomizeTemplate,
			},
		},
	}

	return addOptions(app)
}
