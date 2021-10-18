package karpenter

import (
	"eksdemo/pkg/kubernetes"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
)

func karpenterDefaultProvisioner() *resource.Resource {
	res := &resource.Resource{
		Options: &resource.CommonOptions{
			Name: "karpenter-default-provisioner",
		},

		Manager: &kubernetes.ResourceManager{
			Template: &template.TextTemplate{
				Template: yamlTemplate,
			},
		},
	}
	return res
}

const yamlTemplate = `
apiVersion: karpenter.sh/v1alpha4
kind: Provisioner
metadata:
  name: default
spec:
  provider:
    instanceProfile: KarpenterNodeInstanceProfile-{{ .ClusterName }}
    capacityTypes: [ "spot" ]
    cluster:
      name: {{ .ClusterName }}
      endpoint: {{ .Cluster.Endpoint }}
  ttlSecondsAfterEmpty: 30
`
