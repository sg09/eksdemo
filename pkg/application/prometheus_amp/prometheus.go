package prometheus_amp

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/amp_workspace"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:    https://github.com/prometheus-operator/kube-prometheus/tree/main/docs
// GitHub:  https://github.com/prometheus-operator/kube-prometheus
// Helm:    https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack
// Repo:    https://quay.io/prometheus-operator/prometheus-operator
// Version: Latest is Chart 41.5.1, PromOperator v0.60.1 (as of 10/23/22)
//          But pinning Previous Chart to 34.10.0 due to breaking API Server graphs for k8s < 1.23
//          https://github.com/prometheus-community/helm-charts/issues/2018

func NewApp() *application.Application {
	options, flags := NewOptions()
	options.AmpWorkspaceOptions = &amp_workspace.AmpWorkspaceOptions{
		CommonOptions: resource.CommonOptions{
			Name: "amazon-managed-prometheus",
		},
	}

	app := &application.Application{
		Command: cmd.Command{
			Name:        "prometheus-amp",
			Description: "Prometheus with Amazon Managed Prometheus (AMP)",
			Aliases:     []string{"prom-amp", "promamp"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "prometheus-amp-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: irsaPolicyDocument,
				},
			}),
			amp_workspace.NewResourceWithOptions(options.AmpWorkspaceOptions),
		},

		Installer: &installer.HelmInstaller{
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

const valuesTemplate = `---
fullnameOverride: prometheus-amp
defaultRules:
  rules:
    alertmanager: false
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
  prometheusScrape: false
prometheus-node-exporter:
  fullnameOverride: node-exporter
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

const postRenderKustomize = `---
resources:
- manifest.yaml
patches:
# The Prometheus Operator cretes a kubelet service for monitoring.
# This patch modifies a flag to Prometheus Operator to use a standard
# name for kubelet service that all prometheus installs use.
# This prevents duplicate kubelet services with multiple prometheus installs
# that causes an issue with the prometheus recording rules
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
