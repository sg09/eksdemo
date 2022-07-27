package vpa

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/template"
)

// Docs:    https://github.com/kubernetes/autoscaler/blob/master/vertical-pod-autoscaler/README.md
// GitHub:  https://github.com/kubernetes/autoscaler/tree/master/vertical-pod-autoscaler
// Helm:    https://github.com/FairwindsOps/charts/tree/master/stable/vpa
// Repo:    k8s.gcr.io/autoscaling/vpa-recommender, k8s.gcr.io/autoscaling/vpa-updater
// Version: Latest is chart 1.4.0, VPA 0.11.0 (as of 07/26/22)

func NewApp() *application.Application {
	options, flags := newOptions()

	app := &application.Application{
		Command: cmd.Command{
			Name:        "vpa",
			Description: "Vertical Pod Autoscaler",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "vpa",
			ReleaseName:   "autoscaling-vpa",
			RepositoryURL: "https://charts.fairwinds.com/stable",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options = options
	app.Flags = flags

	return app
}

const valuesTemplate = `---
fullnameOverride: vpa
recommender:
  image:
    tag: {{ .Version }}
updater:
  image:
    tag: {{ .Version }}
admissionController:
  enabled: {{ .AdmissionControllerEnabled }}
  image:
    tag: {{ .Version }}
`
