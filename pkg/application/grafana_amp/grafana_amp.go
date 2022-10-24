package grafana_amp

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:    https://grafana.com/docs/
// GitHub:  https://github.com/grafana/grafana
// Helm:    https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack
// Helm:    https://github.com/grafana/helm-charts/tree/main/charts/grafana
// Repo:    https://hub.docker.com/r/grafana/grafana
// Version: Latest is Chart 41.5.1, Grafana v9.1.7 (as of 10/23/22)
//          But pinning Previous Chart to 34.10.0 due to breaking API Server graphs for k8s < 1.23
//          https://github.com/prometheus-community/helm-charts/issues/2018

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
				PolicyDocTemplate: &template.TextTemplate{
					Template: irsaPolicyDocument,
				},
			}),
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "kube-prometheus-stack",
			ReleaseName:   "grafana-amp",
			RepositoryURL: "https://prometheus-community.github.io/helm-charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
			PostRenderKustomize: &template.TextTemplate{
				Template: postRenderKustomize,
			},
		},
	}
	app.Options, app.Flags = newOptions()

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

const valuesTemplate = `---
fullnameOverride: grafana-amp
defaultRules:
  create: false
alertmanager:
  enabled: false
grafana:
  adminPassword: {{ .GrafanaAdminPassword }}
  fullnameOverride: grafana
  grafana.ini:
    auth:
      sigv4_auth_enabled: true
  image:
    tag: {{ .Version }}
  service:
    type: {{ .ServiceType }}
    annotations:
      {{- .ServiceAnnotations | nindent 6 }}
  # Temporary fix for Issue: https://github.com/prometheus-community/helm-charts/issues/1867
  serviceMonitor:
    labels:
      release: prometheus-amp
{{- if .IngressHost }}
  ingress:
    enabled: true
    ingressClassName: {{ .IngressClass }}
    annotations:
      {{- .IngressAnnotations | nindent 6 }}
    hosts:
    - {{ .IngressHost }}
    tls:
    - hosts:
      - {{ .IngressHost }}
    {{- if ne .IngressClass "alb" }}
      secretName: grafana-amp-tls
    {{- end}}
{{- end }}
  sidecar:
    datasources:
      defaultDatasourceEnabled: false
  additionalDataSources:
  - name: Prometheus
    type: prometheus
    url: {{ .AmpEndpoint }}
    access: proxy
    isDefault: true
    jsonData:
      sigV4Auth: true
      sigV4AuthType: default
      sigV4Region: {{ .Region }}
    timeInterval: 30s
  serviceAccount:
    name: {{ .ServiceAccount }}
    annotations:
      {{ .IrsaAnnotation }}
kubeApiServer:
  # Enabled to create the API Server dashboard, ServiceMonitor deleted in Post Render
  enabled: true
kubelet:
  # Enabled to create the Kubelet dashboard, ServiceMonitor deleted in Post Render
  enabled: true
kubeControllerManager:
  enabled: false
coreDns:
  # Enabled to create the Core DNS dashboard, ServiceMonitor and Service deleted in Post Render
  enabled: true
kubeEtcd:
  enabled: false
kubeScheduler:
  enabled: false
kubeProxy:
  # Enabled to create the kube-proxy dashboard
  enabled: true
  service:
    enabled: false
  serviceMonitor:
    enabled: false
kubeStateMetrics:
  enabled: false
nodeExporter:
  # Enabled to create the USE dashboards, Daemonset and Service deleted in Post Render
  enabled: true
prometheus-node-exporter:
  serviceAccount:
    create: false
  namespaceOverride: delete-me
prometheusOperator:
  enabled: false
prometheus:
  enabled: false
  prometheusSpec:
    remoteWriteDashboards: true
`

// https://github.com/kubernetes-sigs/kustomize/blob/master/examples/patchMultipleObjects.md
const postRenderKustomize = `---
resources:
- manifest.yaml
patches:
# Delete the AlertManager dashboard as it's disabled, alerts will be in AMP
- patch: |-
    $patch: delete
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: needed-but-not-used
  target:
    name: ".*alertmanager-overview$"
# Delete ServiceMonitors that are created by enabling the dashboards for those services
- patch: |-
    $patch: delete
    apiVersion: unused
    kind: unused
    metadata:
      name: unused
  target:
    kind: ServiceMonitor
    name: grafana-amp.*
# Delete the Service that is created by enabling the dashboard for CoreDNS
- patch: |-
    $patch: delete
    apiVersion: unused
    kind: unused
    metadata:
      name: unused
  target:
    kind: Service
    namespace: kube-system
# Delete the Node Exporter Daemonset and Service
- patch: |-
    $patch: delete
    apiVersion: unused
    kind: unused
    metadata:
      name: unused
  target:
    namespace: delete-me
`
