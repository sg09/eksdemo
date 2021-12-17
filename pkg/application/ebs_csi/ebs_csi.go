package ebs_csi

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:    https://github.com/kubernetes-sigs/aws-ebs-csi-driver/tree/master/docs
// GitHub:  https://github.com/kubernetes-sigs/aws-ebs-csi-driver
// Helm:    https://github.com/kubernetes-sigs/aws-ebs-csi-driver/tree/master/charts/aws-ebs-csi-driver
// Repo:    k8s.gcr.io/provider-aws/aws-ebs-csi-driver
// Version: Latest is v1.4.0 (as of 10/20/21)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "ebs-csi",
			Description: "CSI driver for Amazon EBS",
			Aliases:     []string{"ebs"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "ebs-csi-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: policyDocument,
				},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "kube-system",
			ServiceAccount: "ebs-csi-controller-sa",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "v1.4.0",
				Previous: "v1.3.1",
			},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "aws-ebs-csi-driver",
			ReleaseName:   "aws-ebs-csi-driver",
			RepositoryURL: "https://kubernetes-sigs.github.io/aws-ebs-csi-driver",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		PostInstallResources: []*resource.Resource{
			gp3StorageClass(),
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
  - ec2:CreateSnapshot
  - ec2:AttachVolume
  - ec2:DetachVolume
  - ec2:ModifyVolume
  - ec2:DescribeAvailabilityZones
  - ec2:DescribeInstances
  - ec2:DescribeSnapshots
  - ec2:DescribeTags
  - ec2:DescribeVolumes
  - ec2:DescribeVolumesModifications
  Resource: "*"
- Effect: Allow
  Action:
  - ec2:CreateTags
  Resource:
  - arn:aws:ec2:*:*:volume/*
  - arn:aws:ec2:*:*:snapshot/*
  Condition:
    StringEquals:
      ec2:CreateAction:
      - CreateVolume
      - CreateSnapshot
- Effect: Allow
  Action:
  - ec2:DeleteTags
  Resource:
  - arn:aws:ec2:*:*:volume/*
  - arn:aws:ec2:*:*:snapshot/*
- Effect: Allow
  Action:
  - ec2:CreateVolume
  Resource: "*"
  Condition:
    StringLike:
      aws:RequestTag/ebs.csi.aws.com/cluster: 'true'
- Effect: Allow
  Action:
  - ec2:CreateVolume
  Resource: "*"
  Condition:
    StringLike:
      aws:RequestTag/CSIVolumeName: "*"
- Effect: Allow
  Action:
  - ec2:CreateVolume
  Resource: "*"
  Condition:
    StringLike:
      aws:RequestTag/kubernetes.io/cluster/*: owned
- Effect: Allow
  Action:
  - ec2:DeleteVolume
  Resource: "*"
  Condition:
    StringLike:
      ec2:ResourceTag/ebs.csi.aws.com/cluster: 'true'
- Effect: Allow
  Action:
  - ec2:DeleteVolume
  Resource: "*"
  Condition:
    StringLike:
      ec2:ResourceTag/CSIVolumeName: "*"
- Effect: Allow
  Action:
  - ec2:DeleteVolume
  Resource: "*"
  Condition:
    StringLike:
      ec2:ResourceTag/kubernetes.io/cluster/*: owned
- Effect: Allow
  Action:
  - ec2:DeleteSnapshot
  Resource: "*"
  Condition:
    StringLike:
      ec2:ResourceTag/CSIVolumeSnapshotName: "*"
- Effect: Allow
  Action:
  - ec2:DeleteSnapshot
  Resource: "*"
  Condition:
    StringLike:
      ec2:ResourceTag/ebs.csi.aws.com/cluster: 'true'
`
