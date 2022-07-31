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
// Repo:    public.ecr.aws/ebs-csi-driver/aws-ebs-csi-driver
// Version: Latest is Chart 2.9.0, App v1.10.0 (as of 07/31/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "ebs-csi",
			Description: "CSI driver for Amazon EBS",
			Aliases:     []string{"ebscsi", "ebs"},
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

		Installer: &installer.HelmInstaller{
			ChartName:     "aws-ebs-csi-driver",
			ReleaseName:   "storage-ebs-csi",
			RepositoryURL: "https://kubernetes-sigs.github.io/aws-ebs-csi-driver",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options, app.Flags = newOptions()

	return app
}

const valuesTemplate = `---
image:
  tag: {{ .Version }}
controller:
  region: {{ .Region }}
  replicaCount: 1
  serviceAccount:
    name: {{ .ServiceAccount }}
    annotations:
      {{ .IrsaAnnotation }}
storageClasses:
- name: gp3
{{- if .DefaultGp3 }}
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
{{- end }}
  parameters:
    csi.storage.k8s.io/fstype: ext4
    type: gp3
  volumeBindingMode: WaitForFirstConsumer
`
