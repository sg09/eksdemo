package cert_manager

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:    https://cert-manager.io/docs/
// GitHub:  https://github.com/cert-manager/cert-manager
// Helm:    https://github.com/cert-manager/cert-manager/tree/master/deploy/charts/cert-manager
// Repo:    quay.io/jetstack/cert-manager-controller
// Version: Latest is v1.11.0 (as of 01/20/23)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "cert-manager",
			Description: "Cloud Native Certificate Management",
			Aliases:     []string{"certmanager"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "cert-manager-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: policyDocument,
				},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "cert-manager",
			ServiceAccount: "cert-manager",
			DefaultVersion: &application.LatestPrevious{
				// For Helm Chart version: https://artifacthub.io/packages/helm/cert-manager/cert-manager
				LatestChart:   "1.11.0",
				Latest:        "v1.11.0",
				PreviousChart: "1.10.0",
				Previous:      "v1.10.0",
			},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "cert-manager",
			ReleaseName:   "cert-manager",
			RepositoryURL: "https://charts.jetstack.io",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		PostInstallResources: []*resource.Resource{
			clusterIssuer(),
		},
	}
	return app
}

const valuesTemplate = `---
installCRDs: true
replicaCount: 1
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
image:
  tag: {{ .Version }}
`

const policyDocument = `
Version: '2012-10-17'
Statement:
- Effect: Allow
  Action:
  - route53:GetChange
  Resource: arn:aws:route53:::change/*
- Effect: Allow
  Action:
  - route53:ChangeResourceRecordSets
  - route53:ListResourceRecordSets
  Resource: arn:aws:route53:::hostedzone/*
- Effect: Allow
  Action: route53:ListHostedZonesByName
  Resource: "*"
`
