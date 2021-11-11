package amp

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

func NewResourceWithOptions(options *AmpOptions) *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "amp",
			Description: "Amazon Managed Prometheus",
			Args:        []string{"ALIAS"},
		},

		Getter: &Getter{},

		Manager: &Manager{},
	}

	res.Options = options

	return res
}
