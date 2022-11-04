package kms_key

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "kms-key",
			Description: "KMS Key",
			Aliases:     []string{"kms-keys", "kmskeys", "kmskey", "kms"},
			Args:        []string{"ALIAS"},
		},

		Getter: &Getter{},
	}

	res.Options, res.GetFlags = newOptions()

	return res
}
