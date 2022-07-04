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
// Version: Latest is Chart 4.1.4 and App v1.2.1 (as of 07/04/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "nginx",
			Description: "NGINX Ingress Controller",
		},

		Options: &application.ApplicationOptions{
			Namespace:      "ingress-nginx",
			ServiceAccount: "ingress-nginx",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "4.1.4",
				Latest:        "v1.2.1",
				PreviousChart: "4.1.4",
				Previous:      "v1.2.1",
			},
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
	return app
}

const valuesTemplate = `
controller:
  image:
    tag: {{ .Version }}
  service:
    annotations:
      service.beta.kubernetes.io/aws-load-balancer-backend-protocol: tcp
      service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled: "true"
      service.beta.kubernetes.io/aws-load-balancer-type: nlb
    externalTrafficPolicy: Local
serviceAccount:
  name: {{ .ServiceAccount }}
`
