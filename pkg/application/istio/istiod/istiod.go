package istiod

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/template"
)

// Docs:    https://istio.io/latest/docs/
// GitHub:  https://github.com/istio/istio
// Helm:    https://github.com/istio/istio/tree/master/manifests/charts/istio-control/istio-discovery
// Repo:    https://hub.docker.com/r/istio/pilot
// Version: Latest is v1.12.1 (as of 12/11/21)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "istiod",
			Description: "Istio Control Plane",
			Aliases:     []string{"control-plane", "control", "cp"},
		},

		Options: &application.ApplicationOptions{
			Namespace:      "istio-system",
			ServiceAccount: "istiod",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "1.12.1",
				Previous: "1.12.0",
			},
			// Service Account name is hard coded in the Chart
			// https://github.com/istio/istio/blob/master/manifests/charts/istio-control/istio-discovery/templates/serviceaccount.yaml#L10
			DisableServiceAccountFlag: true,
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "istiod",
			ReleaseName:   "istio-istiod",
			RepositoryURL: "https://istio-release.storage.googleapis.com/charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		// TODO: option to choose namespace for monitors

		// PostInstallResources: []*resource.Resource{
		// 	podMonitor(),
		// 	serviceMonitor(),
		// },
	}
	return app
}

const valuesTemplate = `
pilot:
  tag: {{ .Version }}
global:
  istioNamespace: {{ .Namespace }}
`
