package fargate_profile

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/eksctl"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "fargate-profile",
			Description: "EKS Fargate Profile",
			Aliases:     []string{"fargate-profiles", "fargateprofiles", "fargateprofile", "fargate", "fp"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &eksctl.ResourceManager{
			Resource: "fargateprofile",
			Template: &template.TextTemplate{
				Template: eksctl.EksctlHeader + eksctlTemplate,
			},
		},
	}

	res.Options, res.CreateFlags = NewOptions()

	return res
}

const eksctlTemplate = `
fargateProfiles:
- name: {{ .Name }}
  selectors:
{{- range .Namespaces }}
  - namespace: {{ . }}
{{- end }}
`
