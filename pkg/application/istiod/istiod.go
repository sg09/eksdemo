package istiod

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/helm"
	"eksdemo/pkg/template"
)

// Docs:    https://istio.io/latest/docs/
// GitHub:  https://github.com/istio/istio
// Helm:    https://github.com/istio/istio/tree/master/manifests/charts/istio-control/istio-discovery
// Repo:    https://hub.docker.com/r/istio/pilot
// Version: Latest is v1.12.0 (as of 11/29/21)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "istiod",
			Description: "Istio Control Plane",
			Aliases:     []string{"istio"},
		},

		Options: &application.ApplicationOptions{
			Namespace:      "istio-system",
			ServiceAccount: "istiod",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "1.12.0",
				Previous: "1.12.0",
			},
			// Service Account name is hard coded in the Chart
			// https://github.com/istio/istio/blob/master/manifests/charts/istio-control/istio-discovery/templates/serviceaccount.yaml#L10
			DisableServiceAccountFlag: true,
		},

		Installer: &helm.HelmInstaller{
			ChartName:     "istiod",
			ReleaseName:   "istiod",
			RepositoryURL: "https://istio-release.storage.googleapis.com/charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return app
}

const valuesTemplate = `
pilot:
  tag: {{ .Version }}
global:
  istioNamespace: {{ .Namespace }}
`
