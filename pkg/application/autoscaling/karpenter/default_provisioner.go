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

const yamlTemplate = `---
apiVersion: karpenter.sh/v1alpha5
kind: Provisioner
metadata:
  name: default
spec:
  requirements:
    - key: karpenter.sh/capacity-type
      operator: In
      values: ["spot"]
  limits:
    resources:
      cpu: 1000
  providerRef:
    name: default
  ttlSecondsAfterEmpty: 30
---
apiVersion: karpenter.k8s.aws/v1alpha1
kind: AWSNodeTemplate
metadata:
  name: default
spec:
  subnetSelector:
    Name: eksctl-{{ .ClusterName }}-cluster/SubnetPrivate*
  securityGroupSelector:
    aws:eks:cluster-name: {{ .ClusterName }}
`
