package external_dns

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/helm"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:   https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md
// GitHub: https://github.com/kubernetes-sigs/external-dns
// Helm:   https://github.com/kubernetes-sigs/external-dns/tree/master/charts/external-dns
// Repo:   k8s.gcr.io/external-dns/external-dns

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "external-dns",
			Description: "External DNS",
			Aliases:     []string{"edns"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "external-dns-irsa",
				},
				PolicyType: irsa.WellKnown,
				Policy:     []string{"externalDNS"},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "external-dns",
			ServiceAccount: "external-dns",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "v0.8.0",
				Previous: "v0.8.0",
			},
		},

		Installer: &helm.HelmInstaller{
			ChartName:     "external-dns",
			ReleaseName:   "external-dns",
			RepositoryURL: "https://kubernetes-sigs.github.io/external-dns",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return app
}

const valuesTemplate = `
image:
  tag: {{ .Version }}
provider: aws
registry: txt
serviceAccount:
  create: true
  annotations:
    {{ .IrsaAnnotation }}
  name: {{ .ServiceAccount }}
txtOwnerId: {{ .ClusterName }}
`
