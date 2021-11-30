package istio_base

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/helm"
	"eksdemo/pkg/template"
)

// Docs: https://github.com/istio/istio/blob/master/manifests/charts/README.md
// Docs: https://github.com/istio/istio/blob/master/manifests/charts/README-helm3.md
// Helm: https://github.com/istio/istio/tree/master/manifests/charts/base

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "istio-base",
			Description: "Istio Base (includes CRDs)",
		},

		Options: &application.ApplicationOptions{
			Namespace: "istio-system",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "none",
				Previous: "none",
			},
			DisableServiceAccountFlag: true,
		},

		Installer: &helm.HelmInstaller{
			ChartName:     "base",
			ReleaseName:   "istio-base",
			RepositoryURL: "https://istio-release.storage.googleapis.com/charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return app
}

const valuesTemplate = ``
