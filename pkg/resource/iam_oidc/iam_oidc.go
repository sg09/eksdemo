package iam_oidc

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "iam-oidc",
			Description: "IAM OIDC Identity Provider",
			Aliases:     []string{"oidc"},
			Args:        []string{"URL"},
		},

		Getter: &Getter{},
	}

	res.Options = newOptions()

	return res
}
