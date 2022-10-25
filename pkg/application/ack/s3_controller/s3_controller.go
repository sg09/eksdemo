package s3_controller

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
// GitHub:  https://github.com/aws-controllers-k8s/s3-controller
// Helm:    https://github.com/aws-controllers-k8s/s3-controller/tree/main/helm
// Repo:    https://gallery.ecr.aws/aws-controllers-k8s/s3-controller
// Version: Latest is v0.1.5 (as of 10/24/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "ack",
			Name:        "s3-controller",
			Description: "ACK S3 Controller",
			Aliases:     []string{"s3"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "ack-s3-controller-irsa",
				},
				// https://github.com/aws-controllers-k8s/s3-controller/blob/main/config/iam/recommended-policy-arn
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/AmazonS3FullAccess"},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "ack-system",
			ServiceAccount: "ack-s3-controller",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "v0.1.5",
				Latest:        "v0.1.5",
				PreviousChart: "v0.1.2",
				Previous:      "v0.1.2",
			},
		},

		Installer: &installer.HelmInstaller{
			ReleaseName:   "ack-s3-controller",
			RepositoryURL: "oci://public.ecr.aws/aws-controllers-k8s/s3-chart",
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
fullnameOverride: ack-s3-controller
aws:
  region: {{ .Region }}
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
`
