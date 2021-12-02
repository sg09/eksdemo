package efs_csi

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:    https://github.com/kubernetes-sigs/aws-efs-csi-driver/tree/master/docs
// GitHub:  https://github.com/kubernetes-sigs/aws-efs-csi-driver
// Helm:    https://github.com/kubernetes-sigs/aws-efs-csi-driver/tree/master/charts/aws-efs-csi-driver
// Repo:    amazon/aws-efs-csi-driver
// Version: Latest is v1.3.4 (as of 10/21/21)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "efs-csi",
			Description: "CSI driver for Amazon EFS",
			Aliases:     []string{"efs"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "efs-csi-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				Policy:     []string{policyDocument},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "kube-system",
			ServiceAccount: "efs-csi-controller-sa",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "v1.3.4",
				Previous: "v1.3.3",
			},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "aws-efs-csi-driver",
			ReleaseName:   "aws-efs-csi-driver",
			RepositoryURL: "https://kubernetes-sigs.github.io/aws-efs-csi-driver",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return app
}

const valuesTemplate = `
replicaCount: 1
controller:
  serviceAccount:
    annotations:
      {{ .IrsaAnnotation }}
    name: {{ .ServiceAccount }}
image:
  tag: {{ .Version }}
`

const policyDocument = `
Version: '2012-10-17'
Statement:
- Effect: Allow
  Action:
  - elasticfilesystem:DescribeAccessPoints
  - elasticfilesystem:DescribeFileSystems
  Resource: "*"
- Effect: Allow
  Action:
  - elasticfilesystem:CreateAccessPoint
  Resource: "*"
  Condition:
    StringLike:
      aws:RequestTag/efs.csi.aws.com/cluster: 'true'
- Effect: Allow
  Action: elasticfilesystem:DeleteAccessPoint
  Resource: "*"
  Condition:
    StringEquals:
      aws:ResourceTag/efs.csi.aws.com/cluster: 'true'
`
