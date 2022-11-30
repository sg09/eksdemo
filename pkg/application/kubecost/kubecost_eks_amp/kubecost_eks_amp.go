package kubecost_eks_amp

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/amp"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

// Docs:    https://docs.aws.amazon.com/eks/latest/userguide/cost-monitoring.html
// Docs:    https://guide.kubecost.com/
// Helm:    https://github.com/kubecost/cost-analyzer-helm-chart/tree/develop/cost-analyzer
// Values:  https://github.com/kubecost/cost-analyzer-helm-chart/blob/develop/cost-analyzer/values-eks-cost-monitoring.yaml
// Values:  https://github.com/kubecost/cost-analyzer-helm-chart/blob/develop/cost-analyzer/values-amp.yaml
// Repo:    https://gallery.ecr.aws/kubecost/cost-model
// Repo:    https://gallery.ecr.aws/kubecost/frontend
// Repo:    https://gallery.ecr.aws/kubecost/prometheus
// Repo:    https://gallery.ecr.aws/bitnami/configmap-reload
// Version: Latest is Chart/App 1.97.0 (as of 11/27/22)

func NewApp() *application.Application {
	options, _ := newOptions()

	app := &application.Application{
		Command: cmd.Command{
			Parent:      "kubecost",
			Name:        "eks-amp",
			Description: "EKS optimized Kubecost using Amazon Managed Prometheus",
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "kubecost-cost-analyzer-irsa",
				},
				PolicyType: irsa.PolicyARNs,
				Policy: []string{
					"arn:aws:iam::aws:policy/AmazonPrometheusQueryAccess",
					"arn:aws:iam::aws:policy/AmazonPrometheusRemoteWriteAccess",
				},
			}),
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name:           "kubecost-prometheus-server-irsa",
					ServiceAccount: "kubecost-prometheus-server",
				},
				PolicyType: irsa.PolicyARNs,
				Policy: []string{
					"arn:aws:iam::aws:policy/AmazonPrometheusQueryAccess",
					"arn:aws:iam::aws:policy/AmazonPrometheusRemoteWriteAccess",
				},
			}),
			amp.NewResourceWithOptions(options.AmpOptions),
		},

		Options: options,

		Installer: &installer.HelmInstaller{
			ChartName:     "cost-analyzer",
			ReleaseName:   "kubecost-eks",
			RepositoryURL: "oci://public.ecr.aws/kubecost/cost-analyzer",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
			PostRenderKustomize: &template.TextTemplate{
				Template: postRenderKustomize,
			},
		},
	}

	return app
}

const valuesTemplate = `---
fullnameOverride: kubecost-cost-analyzer
global:
  amp:
    enabled: true
    prometheusServerEndpoint: http://localhost:8005/workspaces/{{ .AmpId }}
    remoteWriteService: {{ .AmpEndpoint }}api/v1/remote_write
    sigv4:
      region: {{ .Region }}
  grafana:
    # If false, Grafana will not be installed
    enabled: false
    # If true, the kubecost frontend will route to your grafana through its service endpoint
    proxy: false
sigV4Proxy:
  region: us-west-2
  host: aps-workspaces.{{ .Region }}.amazonaws.com
podSecurityPolicy:
  enabled: false
imageVersion: prod-{{ .Version }}
kubecostFrontend:
  image: public.ecr.aws/kubecost/frontend
kubecostModel:
  image: public.ecr.aws/kubecost/cost-model
  # The total number of days the ETL storage will build
  etlStoreDurationDays: 120
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
{{- if .IngressHost }}
ingress:
  enabled: true
  className: {{ .IngressClass }}
  annotations:
    {{- .IngressAnnotations | nindent 4 }}
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
service:
  type: {{ .ServiceType }}
  annotations: 
    {{- .ServiceAnnotations | nindent 4 }}
prometheus:
  kube-state-metrics:
    fullnameOverride: kubecost-kube-state-metrics
  server:
    fullnameOverride: kubecost-prometheus-server
    image:
      repository: public.ecr.aws/kubecost/prometheus
      tag: v2.35.0
    global:
      # overrides kubecost default of 60s, sets to prom default of 10s
      scrape_timeout: 10s
      external_labels:
        cluster_id: {{ .ClusterName }} # Each cluster should have a unique ID
  configmapReload:
    prometheus:
      image:
        repository: public.ecr.aws/bitnami/configmap-reload
        tag: 0.7.1
  nodeExporter:
    enabled: false
reporting:
  productAnalytics: false
kubecostProductConfigs:
  clusterName: {{ .ClusterName }}
`

const postRenderKustomize = `---
resources:
- manifest.yaml
patches:
# The Prometheus Chart doesn't allow for annotation of the service accounts
# Open issue: https://github.com/kubecost/cost-analyzer-helm-chart/issues/1009
- patch: |-
    - op: add
      path: /metadata/annotations
      value:
        {{ .IrsaAnnotationFor "kubecost-prometheus-server" }}
  target:
    kind: ServiceAccount
    name: kubecost-prometheus-server
`
