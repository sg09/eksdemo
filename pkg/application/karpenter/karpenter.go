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
// Repo:    https://gallery.ecr.aws/karpenter/controller
// Version: Latest is v0.13.1 (as of 07/04/22)

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
				PolicyDocTemplate: &template.TextTemplate{
					Template: irsaPolicyDocument,
				},
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

		PostInstallResources: []*resource.Resource{
			karpenterDefaultProvisioner(),
		},
	}
	app.Options, app.Flags = NewOptions()

	return app
}

const irsaPolicyDocument = `
Version: "2012-10-17"
Statement:
- Effect: Allow
  Resource: "*"
  Action:
  # Write Operations
  - ec2:CreateLaunchTemplate
  - ec2:CreateFleet
  - ec2:RunInstances
  - ec2:CreateTags
  - iam:PassRole
  - ec2:TerminateInstances
  - ec2:DeleteLaunchTemplate
  # Read Operations
  - ec2:DescribeLaunchTemplates
  - ec2:DescribeInstances
  - ec2:DescribeSecurityGroups
  - ec2:DescribeSubnets
  - ec2:DescribeInstanceTypes
  - ec2:DescribeInstanceTypeOfferings
  - ec2:DescribeAvailabilityZones
  - ssm:GetParameter
`

const valuesTemplate = `
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
replicas: 1
controller:
  image: public.ecr.aws/karpenter/controller:{{ .Version }}
clusterName: {{ .ClusterName }}
clusterEndpoint: {{ .Cluster.Endpoint }}
aws:
  defaultInstanceProfile: KarpenterNodeInstanceProfile-{{ .ClusterName }}
`
