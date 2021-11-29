package addon

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/eksctl"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "addon",
			Description: "EKS Managed Addon",
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &eksctl.ResourceManager{
			Resource: "addon",
			Template: &template.TextTemplate{
				Template: eksctl.EksctlHeader + eksctlTemplate,
			},
		},
	}

	res.Options, res.Flags = NewOptions()

	return res
}

const eksctlTemplate = `
addons:
- name: {{ .Name }}
{{- if .Version }}
  version: {{ .Version }}
{{- end }}
`
