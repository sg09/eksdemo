package cluster

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/eksctl"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/nodegroup"
	"eksdemo/pkg/template"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "cluster",
			Description: "EKS Cluster",
		},

		Getter: &Getter{},

		Manager: &eksctl.ResourceManager{
			Resource: "cluster",
			Template: &template.TextTemplate{
				Template: eksctl.EksctlHeader + EksctlTemplate + nodegroup.EksctlTemplate,
			},
		},
	}

	return addOptions(res)
}

const EksctlTemplate = `
addons:
- name: vpc-cni

cloudWatch:
  clusterLogging:
    enableTypes: ["*"]

{{ if .Fargate }}
fargateProfiles:
- name: default
  selectors:
    - namespace: fargate
{{ end }}

{{ if not .NoRoles }}
iam:
  withOIDC: true
  serviceAccounts:
{{- range .IrsaRoles }}
{{- $.IrsaTemplate.Render .Options }}
{{- end }}
{{ end }}
`
