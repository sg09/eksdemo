package karpenter

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/eksctl"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/iam_auth"
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
			iam_auth.NewResourceWithOptions(&iam_auth.IamAuthOptions{
				CommonOptions: resource.CommonOptions{
					Name: "karpenter-node-iam-auth",
				},
				IamAuth: eksctl.IamAuth{
					Arn:      "arn:aws:iam::{{ .Account }}:role/KarpenterNodeRole-{{ .ClusterName }}",
					Groups:   []string{"system:bootstrappers", "system:nodes"},
					Username: "system:node:{{EC2PrivateDNSName}}",
				},
			}),
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "karpenter",
			ReleaseName:   "karpenter",
			RepositoryURL: "https://charts.karpenter.sh",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
			Wait: true,
		},
	}
	app.Options, app.Flags = NewOptions()

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
serviceAccount:
  annotations:
    {{ .IrsaAnnotation }}
  name: {{ .ServiceAccount }}
controller:
  image: public.ecr.aws/karpenter/controller:{{ .Version }}
  clusterName: {{ .ClusterName }}
  clusterEndpoint: {{ .Cluster.Endpoint }}
  replicaCount: 1
`
