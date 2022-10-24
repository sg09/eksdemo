package kube_prometheus

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
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
	app := &application.Application{
		Command: cmd.Command{
			Name:        "kube-prometheus",
			Description: "Kube Prometheus Stack",
			Aliases:     []string{"kube-prom", "kubeprom", "kprom"},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "kube-prometheus-stack",
			ReleaseName:   "kube-prometheus",
			RepositoryURL: "https://prometheus-community.github.io/helm-charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options, app.Flags = newOptions()

	return app
}

const valuesTemplate = `---
fullnameOverride: prometheus
grafana:
  adminPassword: {{ .GrafanaAdminPassword }}
  fullnameOverride: grafana
  # Temporary fix for issue: https://github.com/prometheus-community/helm-charts/issues/1867
  serviceMonitor:
    labels:
      release: kube-prometheus
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
  service:
    annotations:
      {{- .ServiceAnnotations | nindent 6 }}
    type: {{ .ServiceType }}
kubeControllerManager:
  enabled: false
kubeEtcd:
  enabled: false
kubeScheduler:
  enabled: false
kube-state-metrics:
  fullnameOverride: kube-state-metrics
  prometheusScrape: false
prometheus-node-exporter:
  fullnameOverride: node-exporter
  service:
    annotations:
      # Remove with null when https://github.com/helm/helm/issues/9136 is fixed
      prometheus.io/scrape: "false"
prometheusOperator:
  image:
    tag: {{ .Version }}
`
