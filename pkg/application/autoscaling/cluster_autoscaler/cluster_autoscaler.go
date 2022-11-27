package cluster_autoscaler

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:    https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/README.md
// GitHub:  https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler/cloudprovider/aws
// Helm:    https://github.com/kubernetes/autoscaler/tree/master/charts/cluster-autoscaler
// Repo:    k8s.gcr.io/autoscaling/cluster-autoscaler
// Version: Latest for k8s 1.23 is v1.23.0 (as of 11/04/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "autoscaling",
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
				LatestChart: "9.21.0",
				Latest: map[string]string{
					"1.24": "v1.24.0",
					"1.23": "v1.23.0",
					"1.22": "v1.22.3",
					"1.21": "v1.21.3",
				},
				PreviousChart: "9.21.0",
				Previous: map[string]string{
					"1.24": "v1.24.0",
					"1.23": "v1.23.0",
					"1.22": "v1.22.2",
					"1.21": "v1.21.2",
				},
			},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "cluster-autoscaler",
			ReleaseName:   "autoscaling-cluster-autoscaler",
			RepositoryURL: "https://kubernetes.github.io/autoscaler",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return app
}

const valuesTemplate = `---
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
