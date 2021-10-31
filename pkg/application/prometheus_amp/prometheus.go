package prometheus_amp

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/helm"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/amp"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:    https://prometheus.io/docs/introduction/overview/
// GitHub:  https://github.com/prometheus/prometheus
// Helm:    https://github.com/prometheus-community/helm-charts/tree/main/charts/prometheus
// Repo:    quay.io/prometheus/prometheus
// Version: Latest is v2.30.3 (as of 10/28/21)

const ampName = "amp"

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "prometheus-amp",
			Description: "Prometheus with Amazon Managed Prometheus (AMP)",
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "prometheus-amp-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				Policy:     []string{irsaPolicyDocument},
			}),
			amp.NewResourceWithOptions(&amp.AmpOptions{
				CommonOptions: resource.CommonOptions{
					Name: "amazon-managed-prometheus",
				},
				AmpName: ampName,
			}),
		},

		Installer: &helm.HelmInstaller{
			ChartName:     "prometheus",
			ReleaseName:   "prometheus",
			RepositoryURL: "https://prometheus-community.github.io/helm-charts",
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
  - aps:RemoteWrite
  - aps:GetSeries
  - aps:GetLabels
  - aps:GetMetricMetadata
  Resource: "*"
`

const valuesTemplate = `
serviceAccounts:
  server:
    name: {{ .ServiceAccount }}
    annotations:
      {{ .IrsaAnnotation }}
server:
  image:
    tag: {{ .Version }}
  remoteWrite:
  - url: {{ .AmpEndpoint }}api/v1/remote_write
    sigv4:
      region: {{ .Region }}
    queue_config:
      max_samples_per_send: 1000
      max_shards: 200
      capacity: 2500
pushgateway:
  enabled: {{ .PushGateway }}
`
