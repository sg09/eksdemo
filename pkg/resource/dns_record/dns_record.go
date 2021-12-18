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
			Aliases:     []string{"dns-records", "resource-records", "records", "dns"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},
	}

	res.Options, res.Flags = NewOptions()

	return res
}
