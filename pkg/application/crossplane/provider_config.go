package crossplane

import (
	"eksdemo/pkg/kubernetes"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
)

func defaultProviderConfig() *resource.Resource {
	res := &resource.Resource{
		Options: &resource.CommonOptions{
			Name: "default-aws-provider-config",
		},

		Manager: &kubernetes.ResourceManager{
			Template: &template.TextTemplate{
				Template: providerConfigManifest,
			},
		},
	}
	return res
}

const providerConfigManifest = `---
apiVersion: aws.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: InjectedIdentity
`
