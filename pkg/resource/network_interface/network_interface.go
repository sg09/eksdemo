package network_interface

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "network-interface",
			Description: "Elastic Network Interface",
			Aliases:     []string{"network-interfaces", "enis", "eni"},
			Args:        []string{"ID"},
		},

		Getter: &Getter{},
	}

	res.Options, res.Flags = NewOptions()

	return res
}
