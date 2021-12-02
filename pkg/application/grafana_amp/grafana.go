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
// Repo:    grafana/grafana
// Version: Latest is v8.2.3 (as of 11/6/21)

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
fullnameOverride: grafana-amp
defaultRules:
  create: false
alertmanager:
  enabled: false
grafana:
  adminPassword: {{ .GrafanaAdminPassword }}
  fullnameOverride: grafana-amp
  grafana.ini:
    auth:
      sigv4_auth_enabled: true
  image:
    tag: {{ .Version }}
  ingress:
    enabled: {{ not .DisableIngress }}
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
  rbac:
    pspEnabled: false
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
  # Enable to create the API Server dashboard, ServiceMonitor deleted in Post Render
  enabled: true
kubelet:
  # Enable to create the Kubelet dashboard, ServiceMonitor deleted in Post Render
  enabled: true
kubeControllerManager:
  enabled: false
coreDns:
  # Enable to create the Core DNS dashboard, ServiceMonitor and Service deleted in Post Render
  enabled: true
kubeEtcd:
  enabled: false
kubeScheduler:
  enabled: false
kubeProxy:
  # Enable to create the kube-proxy dashboard
  enabled: true
  service:
    enabled: false
  serviceMonitor:
    enabled: false
kubeStateMetrics:
  enabled: false
nodeExporter:
  # Enable to create the USE dashboards, Daemonset and Service deleted in Post Render
  enabled: true
prometheus-node-exporter:
  serviceAccount:
    create: false
  rbac:
    pspEnabled: false
  namespaceOverride: delete-me
prometheusOperator:
  enabled: false
prometheus:
  enabled: false
  prometheusSpec:
    remoteWriteDashboards: true
`

const postRenderKustomize = `
resources:
- manifest.yaml
patches:
# Delete ServiceMonitors that are created when enabling the dashboards for those services
- patch: |-
    $patch: delete
    apiVersion: monitoring.coreos.com/v1
    kind: ServiceMonitor
    metadata:
      name: "*"
  target:
    kind: ServiceMonitor
# Delete the Service that is created when enabling the dashboard for CoreDNS
- patch: |-
    $patch: delete
    apiVersion: v1
    kind: Service
    metadata:
      name: "*"
  target:
    kind: Service
    namespace: kube-system
# Delete the Node Exporter Daemonset
- patch: |-
    $patch: delete
    apiVersion: apps/v1
    kind: DaemonSet
    metadata:
      name: "*"
  target:
    namespace: delete-me
# Delete the Node Exporter Service
- patch: |-
    $patch: delete
    apiVersion: v1
    kind: Service
    metadata:
      name: "*"
  target:
    namespace: delete-me
`
