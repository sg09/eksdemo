package target_group

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "target-group",
			Description: "Target Group",
			Aliases:     []string{"target-groups", "tg"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &Manager{},
	}
	res.Options, res.CreateFlags, res.GetFlags = newOptions()

	return res
}
