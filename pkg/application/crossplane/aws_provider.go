package crossplane

import (
	"eksdemo/pkg/kubernetes"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

type AwsProviderOptions struct {
	irsa.IrsaOptions
	Version *string
}

func awsProvider(o *CrossplaneOptions) *resource.Resource {
	res := &resource.Resource{
		Options: &AwsProviderOptions{
			IrsaOptions: irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					ClusterName:    o.Common().ClusterName,
					Name:           "crossplane-aws-provider",
					Namespace:      o.Common().Namespace,
					ServiceAccount: o.Common().ServiceAccount,
				},
			},
			Version: &o.ProviderVersion,
		},

		Manager: &kubernetes.ResourceManager{
			Template: &template.TextTemplate{
				Template: yamlTemplate,
			},
		},
	}
	return res
}

const yamlTemplate = `---
apiVersion: pkg.crossplane.io/v1alpha1
kind: ControllerConfig
metadata:
  name: aws-config
  annotations:
    {{ .IrsaAnnotation }}
spec:
  podSecurityContext:
    fsGroup: 2000
---
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-aws
spec:
  package: crossplane/provider-aws:{{ .Version }}
  controllerConfigRef:
    name: aws-config
`
