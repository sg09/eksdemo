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

// Docs:    https://github.com/prometheus-operator/kube-prometheus/tree/main/docs
// GitHub:  https://github.com/prometheus-operator/kube-prometheus
// Helm:    https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack
// Repo:    https://quay.io/prometheus-operator/prometheus-operator
// Version: Latest is v0.52.0 (as of 11/7/21)

func NewApp() *application.Application {
	options, flags := NewOptions()
	options.AmpOptions = &amp.AmpOptions{
		CommonOptions: resource.CommonOptions{
			Name: "amazon-managed-prometheus",
		},
	}

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
			amp.NewResourceWithOptions(options.AmpOptions),
		},

		Installer: &helm.HelmInstaller{
			ChartName:     "kube-prometheus-stack",
			ReleaseName:   "prometheus-amp",
			RepositoryURL: "https://prometheus-community.github.io/helm-charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
			PostRenderKustomize: &template.TextTemplate{
				Template: postRenderKustomize,
			},
		},
	}

	app.Options = options
	app.Flags = flags

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
fullnameOverride: prometheus-amp
defaultRules:
  rules:
    alertmanager: false
global:
  rbac:
    pspEnabled: false
alertmanager:
  enabled: false
grafana:
  enabled: false
kubeControllerManager:
  enabled: false
kubeEtcd:
  enabled: false
kubeScheduler:
  enabled: false
kube-state-metrics:
  fullnameOverride: kube-state-metrics-amp
  podSecurityPolicy:
    enabled: false
  prometheusScrape: false
prometheus-node-exporter:
  fullnameOverride: node-exporter
  rbac:
    pspEnabled: false
  service:
    annotations:
      # Remove with null when https://github.com/helm/helm/issues/9136 is fixed
      prometheus.io/scrape: "false"
    # Don't conflict with kube-prometheus install
    port: 9101
    targetPort: 9101
prometheusOperator:
  image:
    tag: {{ .Version }}
prometheus:
  serviceAccount:
    name: {{ .ServiceAccount }}
    annotations:
      {{ .IrsaAnnotation }}
  prometheusSpec:
    remoteWrite:
    - url: {{ .AmpEndpoint }}api/v1/remote_write
      sigv4:
        region: {{ .Region }}
      queueConfig:
        maxSamplesPerSend: 1000
        maxShards: 200
        capacity: 2500
    scrapeInterval: 30s
`

const postRenderKustomize = `
resources:
- manifest.yaml
patches:
# Create a kubelet service that all prometheus installs use, to prevent duplicate 
# kubelet services that cause issues with the prometheus recording rules
- patch: |-
    - op: replace
      path: /spec/template/spec/containers/0/args/0
      value: "--kubelet-service=kube-system/prometheus-kubelet"
  target:
    group: apps
    version: v1
    kind: Deployment
    namespace: {{ .Namespace }}
    name: prometheus-amp-operator
`
