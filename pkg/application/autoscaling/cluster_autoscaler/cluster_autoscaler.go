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
// Version: Latest for k8s 1.24 is v1.24.0 (as of 1/06/23)

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
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: policyDocument,
				},
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

// Policy Notes
//
// autoscaling:DescribeScalingActivities        v1.24 https://github.com/kubernetes/autoscaler/pull/4489
// autoscaling:DescribeTags          removed in v1.25 https://github.com/kubernetes/autoscaler/pull/4424
// ec2:DescribeImages                           v1.26 https://github.com/kubernetes/autoscaler/pull/4588
// ec2:DescribeInstanceTypes                    v1.23 https://github.com/kubernetes/autoscaler/pull/4468
// ec2:GetInstanceTypesFromInstanceRequirements v1.26 https://github.com/kubernetes/autoscaler/pull/4588
// eks:DescribeNodegroup                              https://github.com/kubernetes/autoscaler/pull/4491

const policyDocument = `
Version: '2012-10-17'
Statement:
- Effect: Allow
  Action:
  - autoscaling:DescribeAutoScalingGroups
  - autoscaling:DescribeAutoScalingInstances
  - autoscaling:DescribeLaunchConfigurations
  - autoscaling:DescribeScalingActivities
  - autoscaling:DescribeTags
  - ec2:DescribeImages
  - ec2:DescribeInstanceTypes
  - ec2:DescribeLaunchTemplateVersions
  - ec2:GetInstanceTypesFromInstanceRequirements
  - eks:DescribeNodegroup
  Resource: "*"
- Effect: Allow
  Action:
  - autoscaling:SetDesiredCapacity
  - autoscaling:TerminateInstanceInAutoScalingGroup
  Resource: "*"
  Condition:
    StringEquals:
      aws:ResourceTag/k8s.io/cluster-autoscaler/{{ .ClusterName }}: owned
`

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
