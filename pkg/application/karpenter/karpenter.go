package karpenter

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/helm"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:    https://karpenter.sh/docs/
// GitHub:  https://github.com/awslabs/karpenter
// Helm:    https://github.com/awslabs/karpenter/tree/main/charts/karpenter
// Repo:    public.ecr.aws/karpenter/controller
// Version: Latest is v0.4.0 (as of 10/15/21)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "karpenter",
			Description: "Karpenter Node Autoscaling",
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "karpenter-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				Policy:     []string{irsaPolicyDocument},
			}),
			karpenterNodeRole(),
		},

		Options: &KarpenterOptions{
			application.ApplicationOptions{
				Namespace:      "karpenter",
				ServiceAccount: "karpenter",
				DefaultVersion: &application.LatestPrevious{
					Latest:   "v0.4.0",
					Previous: "v0.3.4",
				},
			},
		},

		Installer: &helm.HelmInstaller{
			ChartName:     "karpenter",
			ReleaseName:   "karpenter",
			RepositoryURL: "https://awslabs.github.io/karpenter/charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return app
}

const irsaPolicyDocument = `
Version: "2012-10-17"
Statement:
- Effect: Allow
  Action:
  # Write Operations
  - ec2:CreateLaunchTemplate
  - ec2:CreateFleet
  - ec2:RunInstances
  - ec2:CreateTags
  - iam:PassRole
  - ec2:TerminateInstances
  # Read Operations
  - ec2:DescribeLaunchTemplates
  - ec2:DescribeInstances
  - ec2:DescribeSecurityGroups
  - ec2:DescribeSubnets
  - ec2:DescribeInstanceTypes
  - ec2:DescribeInstanceTypeOfferings
  - ec2:DescribeAvailabilityZones
  - ssm:GetParameter
  Resource: "*"
`

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
