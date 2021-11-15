package amg

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

func NewResourceWithOptions(options *AmgOptions) *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "amg",
			Description: "Amazon Managed Grafana",
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &Manager{},
	}

	res.Options = options

	return res
}
