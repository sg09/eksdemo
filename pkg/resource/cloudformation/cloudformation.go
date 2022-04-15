package cloudformation

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	options, flags := NewOptions()
	res := NewResourceWithOptions(options)
	res.Flags = flags

	return res
}

func NewResourceWithOptions(options *CloudFormationOptions) *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "cloudformation",
			Description: "CloudFormation Stack",
			Aliases:     []string{"cf"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &Manager{},
	}

	res.Options = options

	return res
}
