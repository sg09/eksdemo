package istio_base

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/template"
)

// Docs:     https://github.com/istio/istio/blob/master/manifests/charts/README.md
// Helm:     https://github.com/istio/istio/tree/master/manifests/charts/base
// Versions: https://artifacthub.io/packages/helm/istio-official/base

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "base",
			Description: "Istio Base (includes CRDs)",
		},

		Options: &application.ApplicationOptions{
			Namespace: "istio-system",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.14.1",
				Latest:        "none",
				PreviousChart: "1.14.0",
				Previous:      "none",
			},
			DisableServiceAccountFlag: true,
		},

		Installer: &installer.HelmInstaller{
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
