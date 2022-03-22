package container_insights_prom

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:     https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/ContainerInsights-Prometheus-Setup.html
// GitHub:   https://github.com/aws-samples/amazon-cloudwatch-container-insights, see Releases for versions
// Manifest: https://github.com/aws-samples/amazon-cloudwatch-container-insights/blob/master/k8s-deployment-manifest-templates/deployment-mode/service/cwagent-prometheus/prometheus-eks.yaml
// Repo:     https://gallery.ecr.aws/cloudwatch-agent/cloudwatch-agent
// Version:  Latest is 1.247350.0b251780 aka k8s/1.3.9 (as of 03/22/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "container-insights-prometheus",
			Description: "CloudWatch Container Insights monitoring for Prometheus",
			Aliases:     []string{"cip", "ciprom"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "container-insights-prometheus-irsa",
				},
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "amazon-cloudwatch",
			ServiceAccount: "cwagent-prometheus",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "1.247350.0b251780",
				Previous: "1.247348.0b251302",
			},
		},

		Installer: &installer.KustomizeInstaller{
			ResourceTemplate: &template.TextTemplate{
				Template: manifestTemplate,
			},
			KustomizeTemplate: &template.TextTemplate{
				Template: kustomizeTemplate,
			},
		},
	}

	return app
}
