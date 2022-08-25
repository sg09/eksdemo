package fsx_lustre_csi

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:    https://github.com/kubernetes-sigs/aws-fsx-csi-driver/tree/master/docs
// GitHub:  https://github.com/kubernetes-sigs/aws-fsx-csi-driver
// Helm:    https://github.com/kubernetes-sigs/aws-fsx-csi-driver/tree/master/charts/aws-fsx-csi-driver
// Repo:    amazon/aws-fsx-csi-driver
// Version: Latest is Chart 1.4.2, App v0.8.2 (as of 07/31/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "storage",
			Name:        "fsx-csi",
			Description: "Amazon FSx for Lustre CSI Driver",
			Aliases:     []string{"fsxcsi", "fsx"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "fsx-csi-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: policyDocument,
				},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "kube-system",
			ServiceAccount: "fsx-csi-controller-sa",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.4.2",
				Latest:        "v0.8.2",
				PreviousChart: "1.4.1",
				Previous:      "v0.8.1",
			},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "aws-fsx-csi-driver",
			ReleaseName:   "storage-fsx-csi",
			RepositoryURL: "https://kubernetes-sigs.github.io/aws-fsx-csi-driver",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return app
}

const valuesTemplate = `---
image:
  tag: {{ .Version }}
controller:
  replicaCount: 1
  serviceAccount:
    name: {{ .ServiceAccount }}
    annotations:
      {{ .IrsaAnnotation }}
`

const policyDocument = `
Version: '2012-10-17'
Statement:
- Effect: Allow
  Action:
  - iam:CreateServiceLinkedRole
  - iam:AttachRolePolicy
  - iam:PutRolePolicy
  Resource: arn:aws:iam::*:role/aws-service-role/s3.data-source.lustre.fsx.amazonaws.com/*
- Action: iam:CreateServiceLinkedRole
  Effect: Allow
  Resource: "*"
  Condition:
    StringLike:
      iam:AWSServiceName:
      - fsx.amazonaws.com
- Effect: Allow
  Action:
  - s3:ListBucket
  - fsx:CreateFileSystem
  - fsx:DeleteFileSystem
  - fsx:DescribeFileSystems
  - fsx:TagResource
  Resource:
  - "*"
`
