package kube_prometheus

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/helm"
	"eksdemo/pkg/template"
)

// Docs:    https://github.com/prometheus-operator/kube-prometheus/tree/main/docs
// GitHub:  https://github.com/prometheus-operator/kube-prometheus
// Helm:    https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack
// Repo:    https://quay.io/prometheus-operator/prometheus-operator
// Version: Latest is v0.52.0 (as of 11/7/21)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "kube-prometheus",
			Description: "Kube Prometheus Stack",
			Aliases:     []string{"kube-prom", "kubeprom", "kprom"},
		},

		Installer: &helm.HelmInstaller{
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
global:
  rbac:
    pspEnabled: false
grafana:
  adminPassword: {{ .GrafanaAdminPassword }}
  fullnameOverride: grafana
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
kubeControllerManager:
  enabled: false
kubeEtcd:
  enabled: false
kubeScheduler:
  enabled: false
kube-state-metrics:
  fullnameOverride: kube-state-metrics
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
prometheusOperator:
  image:
    tag: {{ .Version }}
`
