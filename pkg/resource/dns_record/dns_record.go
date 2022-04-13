package dns_record

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "dns-record",
			Description: "Route53 Resource Record Sets",
			Aliases:     []string{"dns-records", "dns"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &Manager{},
	}

	res.Options, res.DeleteFlags, res.GetFlags = NewOptions()

	return res
}
