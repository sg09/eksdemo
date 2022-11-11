package network_acl_rule

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "network-acl-rule",
			Description: "Network ACL",
			Aliases:     []string{"nacl-rules", "nacl-rule", "naclrules", "naclrule"},
		},

		Getter: &Getter{},
	}

	res.Options, res.GetFlags = NewOptions()

	return res
}
