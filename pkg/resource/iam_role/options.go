package iam_role

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

type IamRoleOptions struct {
	resource.CommonOptions

	All      bool
	LastUsed bool
}

func NewOptions() (options *IamRoleOptions, getFlags cmd.Flags) {
	options = &IamRoleOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	getFlags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "all",
				Description: "show all roles including service roles",
				Shorthand:   "A",
			},
			Option: &options.All,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "last-used",
				Description: "show last used date",
				Shorthand:   "L",
			},
			Option: &options.LastUsed,
		},
	}

	return
}
