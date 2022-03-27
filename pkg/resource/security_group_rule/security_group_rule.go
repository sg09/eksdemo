package security_group_rule

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "security-group-rule",
			Description: "ESecurity Group Rule",
			Aliases:     []string{"security-group-rules", "sg-rules", "sgrules", "sgr"},
			Args:        []string{"ID"},
		},

		Getter: &Getter{},
	}

	res.Options, res.Flags = NewOptions()

	return res
}
