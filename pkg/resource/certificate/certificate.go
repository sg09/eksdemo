package certificate

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "certificate",
			Description: "ACM Cerificate",
			Aliases:     []string{"certificates", "certs", "cert", "acm"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &Manager{},
	}

	res.Options, res.Flags = NewOptions()

	return res
}
