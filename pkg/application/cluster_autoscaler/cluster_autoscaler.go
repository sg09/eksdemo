package cluster_autoscaler

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:   https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/README.md
// GitHub: https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler/cloudprovider/aws
// Helm:   https://github.com/kubernetes/autoscaler/tree/master/charts/cluster-autoscaler
// Repo:   k8s.gcr.io/autoscaling/cluster-autoscaler

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "cluster-autoscaler",
			Description: "Kubernetes Cluster Autoscaler",
			Aliases:     []string{"ca"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "cluster-autoscaler-irsa",
				},
				PolicyType: irsa.WellKnown,
				Policy:     []string{"autoScaler"},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "kube-system",
			ServiceAccount: "cluster-autoscaler",
			DefaultVersion: &application.KubernetesVersionDependent{
				Latest: map[string]string{
					"1.21": "v1.21.1",
					"1.20": "v1.20.1",
					"1.19": "v1.19.2",
					"1.18": "v1.18.3",
					"1.17": "v1.17.4",
				},
				Previous: map[string]string{
					"1.21": "v1.21.0",
					"1.20": "v1.20.0",
					"1.19": "v1.19.1",
					"1.18": "v1.18.2",
					"1.17": "v1.17.3",
				},
			},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "cluster-autoscaler",
			ReleaseName:   "cluster-autoscaler",
			RepositoryURL: "https://kubernetes.github.io/autoscaler",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return app
}

const valuesTemplate = `
autoDiscovery:
  clusterName: {{ .ClusterName }}
awsRegion: {{ .Region }}
cloudProvider: aws
extraArgs:
  balance-similar-node-groups: true
  expander: least-waste
  skip-nodes-with-local-storage: false
  skip-nodes-with-system-pods: false
fullnameOverride: cluster-autoscaler
image:
  tag: {{ .Version }}
rbac:
  create: true
  serviceAccount:
    annotations:
      {{ .IrsaAnnotation }}
    create: true
    name: {{ .ServiceAccount }}
`
