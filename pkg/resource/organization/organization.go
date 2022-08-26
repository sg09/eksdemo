package organization

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	options, flags := NewOptions()
	res := NewResourceWithOptions(options)
	res.CreateFlags = flags

	return res
}

func NewResourceWithOptions(options *OrganizationOptions) *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "organization",
			Description: "AWS Organization",
			Aliases:     []string{"org"},
		},

		Getter: &Getter{},

		Manager: &Manager{},
	}

	res.Options = options

	return res
}
