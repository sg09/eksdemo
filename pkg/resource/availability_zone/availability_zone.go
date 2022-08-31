package availability_zone

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "availability-zone",
			Description: "Availability Zone",
			Aliases:     []string{"availability-zones", "zones", "zone", "az"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},
	}
	res.Options, res.GetFlags = newOptions()

	return res
}
