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
// Version: Latest is v1.7.0 (as of 06/22/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "ebs-csi",
			Description: "CSI driver for Amazon EBS",
			Aliases:     []string{"aws-ebs-csi-driver", "ebscsi", "ebs"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "ebs-csi-irsa",
				},
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy"},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "kube-system",
			ServiceAccount: "ebs-csi-controller-sa",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "2.6.11",
				Latest:        "v1.7.0",
				PreviousChart: "2.6.8",
				Previous:      "v1.6.2",
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

const valuesTemplate = `---
controller:
  replicaCount: 1
  serviceAccount:
    annotations:
      {{ .IrsaAnnotation }}
    name: {{ .ServiceAccount }}
image:
  tag: {{ .Version }}
`
