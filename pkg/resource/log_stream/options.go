package log_stream

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
	"fmt"

	"github.com/spf13/cobra"
)

type LogStreamOptions struct {
	resource.CommonOptions

	LogGroupName string
}

func newOptions() (options *LogStreamOptions, getFlags cmd.Flags) {
	options = &LogStreamOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	clusterFlag := options.NewClusterFlag(resource.Get, false)
	clusterFlag.Validate = nil

	getFlags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "group-name",
				Description: "log group name",
				Shorthand:   "g",
				Validate: func(cmd *cobra.Command, args []string) error {
					if options.LogGroupName == "" && options.ClusterName == "" {
						return fmt.Errorf("must include either %q or %q flag", "--group-name", "--cluster")
					}

					if options.LogGroupName == "" {
						options.LogGroupName = fmt.Sprintf("/aws/eks/%s/cluster", options.ClusterName)
					}

					return nil
				},
			},
			Option: &options.LogGroupName,
		},
		clusterFlag,
	}

	return
}
