package argo_cd

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/template"
)

// Docs:    https://argo-cd.readthedocs.io/
// GitHub:  https://github.com/argoproj/argo-cd
// Helm:    https://github.com/argoproj/argo-helm/tree/main/charts/argo-cd
// Repo:    quay.io/argoproj/argocd
// Version: Latest Chart is 4.9.14, Argo CD v2.4.6 (as of 07/14/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "cd",
			Description: "Declarative continuous deployment for Kubernetes",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "argo-cd",
			ReleaseName:   "argo-cd",
			RepositoryURL: "https://argoproj.github.io/argo-helm",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options, app.Flags = newOptions()

	return app
}

const valuesTemplate = `
fullnameOverride: argocd
global:
  image:
    tag: {{ .Version }}
server:
{{- if .IngressHost }}
  extraArgs:
    - --insecure
{{- end }}
  certificate:
    # -- Deploy a Certificate resource (requires cert-manager)
    enabled: false
{{- if not .IngressHost }}
  service:
    type: LoadBalancer
{{- else }}
  ingress:
    enabled: true
    annotations:
    {{- if eq .IngressClass "alb" }}
      alb.ingress.kubernetes.io/listen-ports: '[{"HTTPS":443}]'
      alb.ingress.kubernetes.io/scheme: internet-facing
      alb.ingress.kubernetes.io/target-type: ip
    {{- end }}
    ingressClassName: {{ .IngressClass }}
    hosts:
      - {{ .IngressHost }}
    tls:
      - hosts:
          - {{ .IngressHost }}
  ingressGrpc:
    enabled: true
  {{- if eq .IngressClass "alb" }}
    isAWSALB: true
    awsALB:
      serviceType: ClusterIP
      backendProtocolVersion: GRPC
  {{- end }}
{{- end }}
configs:
  secret:
    # -- Bcrypt hashed admin password
    argocdServerAdminPassword: "{{ .AdminPassword | bcrypt }}"
`
