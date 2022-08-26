package nodegroup

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewGravitonResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "nodegroup-graviton",
			Description: "Managed Nodegroup with Graviton Instances",
			Aliases:     []string{"graviton", "ng-graviton"},
			Args:        []string{"NAME"},
		},
	}

	res.Options = &resource.CommonOptions{}
	res.CreateFlags = cmd.Flags{}

	return res
}
