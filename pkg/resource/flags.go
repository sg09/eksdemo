package resource

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cmd"
	"fmt"
)

func (o *CommonOptions) NewClusterFlag(action Action) *cmd.StringFlag {
	flag := &cmd.StringFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "cluster",
			Description: fmt.Sprintf("cluster to %s resource (required)", action),
			Shorthand:   "c",
			Required:    true,
			Validate: func() error {
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

func (o *CommonOptions) NewNamespaceFlag(action Action) *cmd.StringFlag {
	flag := &cmd.StringFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "namespace",
			Description: fmt.Sprintf("namespace to %s resource (required)", action),
			Shorthand:   "n",
			Required:    true,
		},
		Option: &o.Namespace,
	}
	return flag
}
