package iam_policy

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "iam-policy",
			Description: "IAM Policy",
			Aliases:     []string{"iam-policies", "policies", "policy"},
			Args:        []string{"POLICY_ARN"},
		},

		Getter: &Getter{},
	}

	res.Options, res.GetFlags = NewOptions()

	return res
}
