package eks_controller

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:    https://aws-controllers-k8s.github.io/community/docs/community/overview/
// Docs:    https://aws-controllers-k8s.github.io/community/reference/
// GitHub:  https://github.com/aws-controllers-k8s/eks-controller
// Helm:    https://github.com/aws-controllers-k8s/eks-controller/tree/main/helm
// Chart:   https://gallery.ecr.aws/aws-controllers-k8s/eks-chart
// Repo:    https://gallery.ecr.aws/aws-controllers-k8s/eks-controller
// Version: Latest is v0.1.2 (as of 06/21/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "eks-controller",
			Description: "ACK EKS Controller",
			Aliases:     []string{"eks"},
		},

		Dependencies: []*resource.Resource{
			fargatePodExecutionRole(),
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "ack-eks-controller-irsa",
				},
				// https://github.com/aws-controllers-k8s/eks-controller/blob/main/config/iam/recommended-inline-policy
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: policyDocTemplate,
				},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "ack-system",
			ServiceAccount: "ack-eks-controller",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "v0.1.2",
				Latest:        "v0.1.2",
				PreviousChart: "v0.1.1",
				Previous:      "v0.1.1",
			},
		},

		Installer: &installer.HelmInstaller{
			ReleaseName:   "ack-eks-controller",
			RepositoryURL: "oci://public.ecr.aws/aws-controllers-k8s/eks-chart",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return app
}

const policyDocTemplate = `
Version: '2012-10-17'
Statement:
- Effect: Allow
  Action:
  - eks:*
  Resource: "*"
- Effect: Allow
  Action:
  - iam:GetRole
  - iam:PassRole
  Resource: "*"
`

// TODO: Can iam:PassRole be restricted? Something like below...
// Resource: arn:aws:iam::{{ .Account }}:role/eksdemo.{{ .ClusterName }}.fargate-pod-execution-role
// Condition:
//   StringLike:
//     "iam:PassedToService": eks-fargate-pods.amazonaws.com

const valuesTemplate = `---
image:
  tag: {{ .Version }}
fullnameOverride: ack-eks-controller
aws:
  region: {{ .Region }}
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
`
