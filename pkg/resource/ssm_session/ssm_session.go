package ssm_session

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "ssm-session",
			Description: "SSM Session",
			Aliases:     []string{"session"},
			Args:        []string{"ID"},
		},

		Getter: &Getter{},

		Manager: &Manager{},
	}
	res.Options, res.GetFlags = newOptions()

	return res
}
