package grafana_amp

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/helm"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:    https://grafana.com/docs/
// GitHub:  https://github.com/grafana/grafana
// Helm:    https://github.com/grafana/helm-charts/tree/main/charts/grafana
// Repo:    grafana/grafana
// Version: Latest is v8.2.2 (as of 10/30/21)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "grafana-amp",
			Description: "Grafana with Amazon Managed Prometheus (AMP)",
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "grafana-amp-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				Policy:     []string{irsaPolicyDocument},
			}),
		},

		Installer: &helm.HelmInstaller{
			ChartName:     "grafana",
			ReleaseName:   "grafana",
			RepositoryURL: "https://grafana.github.io/helm-charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options, app.Flags = NewOptions()

	return app
}

const irsaPolicyDocument = `
Version: "2012-10-17"
Statement:
- Effect: Allow
  Action:
  - aps:QueryMetrics
  - aps:GetSeries
  - aps:GetLabels
  - aps:GetMetricMetadata
  Resource: "*"
`

const valuesTemplate = `
adminPassword: {{ .GrafanaAdminPassword }}
datasources:
  datasources.yaml:
    apiVersion: 1
    datasources:
    - name: Prometheus
      type: prometheus
      url: {{ .AmpEndpoint }}
      access: proxy
      isDefault: true
      jsonData:
        sigV4Auth: true
        sigV4AuthType: default
        sigV4Region: {{ .Region }}
image:
  tag: {{ .Version }}
ingress:
  enabled: {{ not .DisableIngress }}
  hosts: []
  annotations:
    kubernetes.io/ingress.class: alb
    alb.ingress.kubernetes.io/scheme: internet-facing
    alb.ingress.kubernetes.io/target-type: 'ip'
  {{- if .TLSHost }}
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTPS":443}]'
  tls:
  - hosts:
    - {{ .TLSHost }}
  {{- end }}
grafana.ini:
  auth:
    sigv4_auth_enabled: true
rbac:
  pspEnabled: false
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
`
