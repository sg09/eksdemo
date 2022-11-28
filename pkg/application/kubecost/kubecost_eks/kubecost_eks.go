package kubecost_eks

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/template"
)

// Docs:    https://docs.aws.amazon.com/eks/latest/userguide/cost-monitoring.html
// Docs:    https://guide.kubecost.com/
// Helm:    https://github.com/kubecost/cost-analyzer-helm-chart/tree/develop/cost-analyzer
// Values:  https://github.com/kubecost/cost-analyzer-helm-chart/blob/develop/cost-analyzer/values-eks-cost-monitoring.yaml
// Repo:    https://gallery.ecr.aws/kubecost/cost-model
// Repo:    https://gallery.ecr.aws/kubecost/frontend
// Repo:    https://gallery.ecr.aws/kubecost/prometheus
// Repo:    https://gallery.ecr.aws/bitnami/configmap-reload
// Version: Latest is Chart/App 1.97.0 (as of 11/20/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "kubecost",
			Name:        "eks",
			Description: "EKS optimized bundle of Kubecost",
		},

		Options: &application.ApplicationOptions{
			ExposeIngressAndLoadBalancer: true,
			Namespace:                    "kubecost",
			ServiceAccount:               "kubecost-cost-analyzer",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.97.0",
				Latest:        "1.97.0",
				PreviousChart: "1.97.0",
				Previous:      "1.97.0",
			},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "cost-analyzer",
			ReleaseName:   "kubecost-eks",
			RepositoryURL: "oci://public.ecr.aws/kubecost/cost-analyzer",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}

	return app
}

const valuesTemplate = `---
fullnameOverride: kubecost-cost-analyzer
global:
  grafana:
    # If false, Grafana will not be installed
    enabled: false
    # If true, the kubecost frontend will route to your grafana through its service endpoint
    proxy: false
podSecurityPolicy:
  enabled: false
imageVersion: prod-{{ .Version }}
kubecostFrontend:
  image: public.ecr.aws/kubecost/frontend
kubecostModel:
  image: public.ecr.aws/kubecost/cost-model
  # The total number of days the ETL storage will build
  etlStoreDurationDays: 120
serviceAccount:
  name: {{ .ServiceAccount }}
{{- if .IngressHost }}
ingress:
  enabled: true
  className: {{ .IngressClass }}
  annotations:
    {{- .IngressAnnotations | nindent 4 }}
  pathType: Prefix
  hosts:
    - {{ .IngressHost }}
  tls:
  - hosts:
    - {{ .IngressHost }}
  {{- if ne .IngressClass "alb" }}
    secretName: cost-analyzer-tls
  {{- end }}
{{- end }}
service:
  type: {{ .ServiceType }}
  annotations: 
    {{- .ServiceAnnotations | nindent 4 }}
prometheus:
  kube-state-metrics:
    fullnameOverride: kubecost-kube-state-metrics
  server:
    fullnameOverride: kubecost-prometheus-server
    image:
      repository: public.ecr.aws/kubecost/prometheus
      tag: v2.35.0
    global:
      # overrides kubecost default of 60s, sets to prom default of 10s
      scrape_timeout: 10s
      external_labels:
        cluster_id: {{ .ClusterName }} # Each cluster should have a unique ID
  configmapReload:
    prometheus:
      image:
        repository: public.ecr.aws/bitnami/configmap-reload
        tag: 0.7.1
  nodeExporter:
    enabled: false
reporting:
  productAnalytics: false
kubecostProductConfigs:
  clusterName: {{ .ClusterName }}
`
