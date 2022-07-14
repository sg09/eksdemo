package wordpress

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/template"
)

// Docs:    https://codex.wordpress.org/Main_Page
// GitHub:  https://github.com/WordPress/WordPress
// Helm:    https://github.com/bitnami/charts/tree/master/bitnami/wordpress
// Repo:    https://hub.docker.com/r/bitnami/wordpress
// Version: Latest is 6.0.0 (as of 06/25/22)

func NewApp() *application.Application {
	options, flags := NewOptions()

	app := &application.Application{
		Command: cmd.Command{
			Name:        "wordpress",
			Description: "WordPress Blog",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "wordpress",
			ReleaseName:   wordpressReleaseName,
			RepositoryURL: "https://charts.bitnami.com/bitnami",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
			PVCLabels: map[string]string{
				"app.kubernetes.io/instance": wordpressReleaseName,
			},
		},
	}

	app.Options = options
	app.Flags = flags

	return app
}

const wordpressReleaseName = `wordpress`

const valuesTemplate = `---
global:
  storageClass: {{ .StorageClass }}
image:
  tag: {{ .Version }}
wordpressPassword: {{ .WordpressPassword }}
{{- if .IngressHost }}
service:
  type: ClusterIP
ingress:
  enabled: true
  pathType: Prefix
  ingressClassName: alb
  hostname: {{ .IngressHost }}
  annotations:
    alb.ingress.kubernetes.io/scheme: internet-facing
    alb.ingress.kubernetes.io/target-type: 'ip'
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTPS":443}]'
  tls: true
{{- end }}
`
