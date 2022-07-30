package kubecost

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/template"
)

// Docs:    https://guide.kubecost.com/
// Helm:    https://github.com/kubecost/cost-analyzer-helm-chart/tree/develop/cost-analyzer
// Repo:    gcr.io/kubecost1/cost-model
// Version: Latest is Chart/App 1.95.0 (as of 07/20/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "kubecost",
			Description: "Monitor & Reduce Kubernetes Spend",
		},

		Options: &application.ApplicationOptions{
			EnableIngress:  true,
			Namespace:      "kubecost",
			ServiceAccount: "kubecost-cost-analyzer",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.95.0",
				Latest:        "1.95.0",
				PreviousChart: "1.94.3",
				Previous:      "1.94.3",
			},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "cost-analyzer",
			ReleaseName:   "kubecost",
			RepositoryURL: "https://kubecost.github.io/cost-analyzer/",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}

	return app
}

const valuesTemplate = `---
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
podSecurityPolicy:
  enabled: false
service:
  type: {{ .ServiceType }}
{{- if eq .ServiceType "LoadBalancer" }}
  annotations: 
    {{- .ServiceAnnotations | nindent 4 }}
{{- end }}
grafana:
  rbac:
    pspEnabled: false
serviceAccount:
  name: {{ .ServiceAccount }}
kubecostProductConfigs:
  clusterName: {{ .ClusterName }}
`
