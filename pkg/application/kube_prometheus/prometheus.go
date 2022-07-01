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
// Version: Latest is Chart 36.2.0, PromOperator v0.57.0 (as of 06/28/22)
//          But pinning to Chart 34.10.0, PromOperator v0.55.0 due to breaking API Server graphs
//          https://github.com/prometheus-community/helm-charts/issues/2018

// TODO: consider no version flag, perhaps mark as "n/a"
//       and instead consider --grafana-version, --kube-state-metrics-version, --prom-operator-version, etc.

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
	return addOptions(app)
}

const valuesTemplate = `
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
    ingressClassName: alb
    annotations:
      alb.ingress.kubernetes.io/scheme: internet-facing
      alb.ingress.kubernetes.io/target-type: 'ip'
      alb.ingress.kubernetes.io/listen-ports: '[{"HTTPS":443}]'
    tls:
    - hosts:
      - {{ .IngressHost }}
{{- else }}
  service:
    type: LoadBalancer
{{- end }}
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
