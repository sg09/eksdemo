package nginx

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/template"
)

// Docs:    https://kubernetes.github.io/ingress-nginx/
// GitHub:  https://github.com/kubernetes/ingress-nginx
// Helm:    https://github.com/kubernetes/ingress-nginx/tree/main/charts/ingress-nginx
// Repo:    registry.k8s.io/ingress-nginx/controller
// Version: Latest is Chart 4.2.0 and App v1.3.0 (as of 07/20/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "nginx",
			Description: "NGINX Ingress Controller",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "ingress-nginx",
			ReleaseName:   "ingress-nginx",
			RepositoryURL: "https://kubernetes.github.io/ingress-nginx",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options, app.Flags = newOptions()

	return app
}

const valuesTemplate = `---
controller:
  image:
    tag: {{ .Version }}
  replicaCount: {{ .Replicas }}
  service:
    annotations:
      service.beta.kubernetes.io/aws-load-balancer-backend-protocol: tcp
      service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled: "true"
      service.beta.kubernetes.io/aws-load-balancer-type: nlb
    externalTrafficPolicy: Local
serviceAccount:
  name: {{ .ServiceAccount }}
`
