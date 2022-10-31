package iam_role

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
	"fmt"

	"github.com/spf13/cobra"
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

	clusterFlag := options.NewClusterFlag(resource.Get, false)
	clusterFlag.Description = "filter by IRSA roles for cluster"
	origValidate := clusterFlag.Validate
	clusterFlag.Validate = func(cmd *cobra.Command, args []string) error {
		if options.ClusterName != "" && len(args) > 0 {
			return fmt.Errorf("%q flag cannot be used with NAME argument", "--cluster")
		}
		return origValidate(cmd, args)
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
		clusterFlag,
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
