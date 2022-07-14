package adot_operator

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/template"
)

// Docs:    https://opentelemetry.io/docs/
// Docs:    https://github.com/open-telemetry/opentelemetry-operator/blob/main/docs/api.md
// GitHub:  https://github.com/open-telemetry/opentelemetry-operator
// Helm:    https://github.com/open-telemetry/opentelemetry-helm-charts
// Repo:    ghcr.io/open-telemetry/opentelemetry-operator/opentelemetry-operator
// Version: Latest is v0.53.0 (as of 06/22/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "adot-operator",
			Description: "AWS Distro for OpenTelemetry Operator",
			Aliases:     []string{"adot", "otel-operator", "otel"},
		},

		Options: &application.ApplicationOptions{
			Namespace: "adot-system",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "0.8.2",
				Latest:        "v0.53.0",
				PreviousChart: "0.7.0",
				Previous:      "v0.51.0",
			},
			// Service Account name isn't flexible
			// https://github.com/open-telemetry/opentelemetry-helm-charts/blob/main/charts/opentelemetry-operator/templates/serviceaccount.yaml
			DisableServiceAccountFlag: true,
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "opentelemetry-operator",
			ReleaseName:   "adot-operator",
			RepositoryURL: "https://open-telemetry.github.io/opentelemetry-helm-charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return app
}

const valuesTemplate = `---
replicaCount: 1
nameOverride: adot-operator
manager:
  image:
    tag: {{ .Version }}
`
