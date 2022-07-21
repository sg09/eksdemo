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

		Installer: &installer.HelmInstaller{
			ChartName:     "cost-analyzer",
			ReleaseName:   "kubecost",
			RepositoryURL: "https://kubecost.github.io/cost-analyzer/",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options, app.Flags = newOptions()

	return app
}

const valuesTemplate = `---
{{- if .IngressHost }}
ingress:
  enabled: true
  className: {{ .IngressClass }}
  annotations:
    {{- .IngressAnnotations | nindent 4 }}
  {{- if .AdminPassword }}
    nginx.ingress.kubernetes.io/auth-type: basic
    nginx.ingress.kubernetes.io/auth-secret: basic-auth
    nginx.ingress.kubernetes.io/auth-realm: "Authentication Required"
  {{- end }}
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
