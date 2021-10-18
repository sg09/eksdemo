package aws_lb

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/helm"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:    https://kubernetes-sigs.github.io/aws-load-balancer-controller/
// GitHub:  https://github.com/kubernetes-sigs/aws-load-balancer-controller
// Helm:    https://github.com/aws/eks-charts/tree/master/stable/aws-load-balancer-controller
// Repo:    602401143452.dkr.ecr.us-west-2.amazonaws.com/amazon/aws-load-balancer-controller
// Version: Latest is v2.2.4 (as of 9/23/21)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "aws-lb",
			Description: "AWS Load Balancer Controller",
			Aliases:     []string{"awslb"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "aws-load-balancer-controller-irsa",
				},
				PolicyType: irsa.WellKnown,
				Policy:     []string{"awsLoadBalancerController"},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "aws-lb",
			ServiceAccount: "aws-lb-controller",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "v2.2.4",
				Previous: "v2.2.3",
			},
		},

		Installer: &helm.HelmInstaller{
			ChartName:     "aws-load-balancer-controller",
			ReleaseName:   "aws-load-balancer-controller",
			RepositoryURL: "https://aws.github.io/eks-charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return app
}

const valuesTemplate = `
clusterName: {{ .ClusterName }}
replicaCount: 1
serviceAccount:
  annotations:
    {{ .IrsaAnnotation }}
  name: {{ .ServiceAccount }}
image:
  tag: {{ .Version }}
`
