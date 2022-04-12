package iam_role

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "iam-role",
			Description: "IAM Role",
			Aliases:     []string{"iam-roles", "iamroles", "iamrole", "roles", "role"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},
	}

	res.Options, res.GetFlags = NewOptions()

	return res
}
