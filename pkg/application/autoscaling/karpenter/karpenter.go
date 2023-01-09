package karpenter

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/eksctl"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/iam_auth"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/resource/service_linked_role"
	"eksdemo/pkg/template"
)

// Docs:    https://karpenter.sh/docs/
// GitHub:  https://github.com/awslabs/karpenter
// Helm:    https://github.com/awslabs/karpenter/tree/main/charts/karpenter
// Repo:    https://gallery.ecr.aws/karpenter/controller
// Version: Latest is v0.21.1 (as of 12/31/22)

func NewApp() *application.Application {
	options, flags := newOptions()

	app := &application.Application{
		Command: cmd.Command{
			Parent:      "autoscaling",
			Name:        "karpenter",
			Description: "Karpenter Node Autoscaling",
		},

		Dependencies: []*resource.Resource{
			service_linked_role.NewResourceWithOptions(&service_linked_role.ServiceLinkedRoleOptions{
				CommonOptions: resource.CommonOptions{
					Name: "ec2-spot-service-linked-role",
				},
				RoleName:    "AWSServiceRoleForEC2Spot",
				ServiceName: "spot.amazonaws.com",
			}),
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
			karpenterSqsQueue(),
			iam_auth.NewResourceWithOptions(&iam_auth.IamAuthOptions{
				CommonOptions: resource.CommonOptions{
					Name: "karpenter-node-iam-auth",
				},
				IamAuth: eksctl.IamAuth{
					Arn:      "arn:{{ .Partition }}:iam::{{ .Account }}:role/KarpenterNodeRole-{{ .ClusterName }}",
					Groups:   []string{"system:bootstrappers", "system:nodes"},
					Username: "system:node:{{EC2PrivateDNSName}}",
				},
			}),
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "karpenter",
			ReleaseName:   "autoscaling-karpenter",
			RepositoryURL: "oci://public.ecr.aws/karpenter/karpenter",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
			Wait: true,
		},

		PostInstallResources: []*resource.Resource{
			karpenterDefaultProvisioner(options),
		},
	}
	app.Options = options
	app.Flags = flags

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
  - ec2:TerminateInstances
  - ec2:DeleteLaunchTemplate
  # Read Operations
  - ec2:DescribeLaunchTemplates
  - ec2:DescribeInstances
  - ec2:DescribeSecurityGroups
  - ec2:DescribeSubnets
  - ec2:DescribeImages
  - ec2:DescribeInstanceTypes
  - ec2:DescribeInstanceTypeOfferings
  - ec2:DescribeAvailabilityZones
  - ec2:DescribeSpotPriceHistory
  - ssm:GetParameter
  - pricing:GetProducts
- Effect: Allow
  Action:
  # Write Operations
  - sqs:DeleteMessage
  # Read Operations
  - sqs:GetQueueUrl
  - sqs:GetQueueAttributes
  - sqs:ReceiveMessage
  Resource: arn:{{ .Partition }}:sqs:{{ .Region }}:{{ .Account }}:karpenter-{{ .ClusterName }}
- Effect: Allow
  Action:
  - iam:PassRole
  Resource: arn:{{ .Partition }}:iam::{{ .Account }}:role/KarpenterNodeRole-{{ .ClusterName }}
`

const valuesTemplate = `---
fullnameOverride: karpenter
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
replicas: 1
controller:
  image: public.ecr.aws/karpenter/controller:{{ .Version }}
settings:
  aws:
    clusterName: {{ .ClusterName }}
    clusterEndpoint: {{ .Cluster.Endpoint }}
    defaultInstanceProfile: KarpenterNodeInstanceProfile-{{ .ClusterName }}
    interruptionQueueName: karpenter-{{ .ClusterName }}
`
