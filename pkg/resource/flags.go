package resource

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cmd"
	"fmt"

	"github.com/spf13/cobra"
)

func (o *CommonOptions) NewClusterFlag(action Action, required bool) *cmd.StringFlag {
	flag := &cmd.StringFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "cluster",
			Description: fmt.Sprintf("cluster to %s resource", action),
			Shorthand:   "c",
			Required:    required,
			Validate: func(cmd *cobra.Command, args []string) error {
				if !required && o.ClusterName == "" {
					return nil
				}

				cluster, err := aws.EksDescribeCluster(o.ClusterName)
				if err != nil {
					return err
				}
				o.Cluster = cluster
				o.KubernetesVersion = aws.StringValue(cluster.Version)

				o.Account = aws.AccountId()
				o.Region = aws.Region()

				return nil
			},
		},
		Option: &o.ClusterName,
	}
	return flag
}

func (o *CommonOptions) NewDryRunFlag() *cmd.BoolFlag {
	flag := &cmd.BoolFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "dry-run",
			Description: "don't create, just print out all creation steps",
		},
		Option: &o.DryRun,
	}
	return flag
}

func (o *CommonOptions) NewIdFlag() *cmd.StringFlag {
	flag := &cmd.StringFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "id",
			Description: "delete by ID instead",
		},
		Option: &o.Id,
	}
	return flag
}

func (o *CommonOptions) NewNamespaceFlag(action Action) *cmd.StringFlag {
	flag := &cmd.StringFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "namespace",
			Description: fmt.Sprintf("namespace to %s resource", action),
			Shorthand:   "n",
		},
		Option: &o.Namespace,
	}
	return flag
}
