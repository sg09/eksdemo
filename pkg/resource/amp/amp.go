package amp

import (
	"eksdemo/pkg/cloudformation"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
)

func NewResource() *resource.Resource {
	options, flags := NewOptions()
	res := NewResourceWithOptions(options)
	res.Flags = flags

	return res
}

func NewResourceWithOptions(options *AmpOptions) *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "amp",
			Description: "Amazon Managed Prometheus",
			Args:        []string{"ALIAS"},
		},

		Getter: &Getter{},

		Manager: &cloudformation.ResourceManager{
			Resource: "amp",
			Template: &template.TextTemplate{
				Template: cloudformationTemplate,
			},
		},
	}

	res.Options = options

	return res
}

const cloudformationTemplate = `
Resources:
  APSWorkspace:
    Type: AWS::APS::Workspace
    Properties:
      Alias: {{ .Alias }}
`
