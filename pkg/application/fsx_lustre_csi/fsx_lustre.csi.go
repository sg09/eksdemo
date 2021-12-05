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
// Version: Latest is v0.5.0 (as of 10/21/21)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "fsx-lustre-csi",
			Description: "CSI Driver of Amazon FSx for Lustre",
			Aliases:     []string{"fsx", "fsx-lustre"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "fsx-lustre-csi-irsa",
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
				Latest:   "v0.5.0",
				Previous: "v0.4.0",
			},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "aws-fsx-csi-driver",
			ReleaseName:   "aws-fsx-csi-driver",
			RepositoryURL: "https://kubernetes-sigs.github.io/aws-fsx-csi-driver",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return app
}

const valuesTemplate = `
controller:
  replicaCount: 1
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
