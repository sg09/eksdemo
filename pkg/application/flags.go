package application

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/kubernetes"
	"fmt"
)

func (o *ApplicationOptions) NewClusterFlag(action Action) *cmd.StringFlag {
	flag := &cmd.StringFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "cluster",
			Description: fmt.Sprintf("cluster to %s application (required)", action),
			Shorthand:   "c",
			Required:    true,
			Validate: func() error {
				cluster, err := aws.EksDescribeCluster(o.ClusterName)
				if err != nil {
					return err
				}

				o.kubeContext, err = kubernetes.KubeContextForCluster(cluster)
				if err != nil {
					return err
				}
				if o.kubeContext == "" {
					return fmt.Errorf("cluster \"%s\" not found in Kubeconfig", o.ClusterName)
				}

				o.Cluster = cluster
				o.Account = aws.AccountId()
				o.Region = aws.Region()

				return nil
			},
		},
		Option: &o.ClusterName,
	}
	return flag
}

func (o *ApplicationOptions) NewDeleteRoleFlag() *cmd.BoolFlag {
	flag := &cmd.BoolFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "delete-dependencies",
			Description: "delete application dependencies",
			Shorthand:   "D",
		},
		Option: &o.DeleteDependencies,
	}
	return flag
}

func (o *ApplicationOptions) NewDryRunFlag() *cmd.BoolFlag {
	flag := &cmd.BoolFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "dry-run",
			Description: "don't install, just print out all installation steps",
		},
		Option: &o.DryRun,
	}
	return flag
}

func (o *ApplicationOptions) NewNamespaceFlag(action Action) *cmd.StringFlag {
	var desc string

	switch action {
	case Install:
		desc = "namespace to install"
	case Uninstall:
		desc = "namespace application is installed"

	}
	flag := &cmd.StringFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "namespace",
			Description: desc,
			Shorthand:   "n",
		},
		Option: &o.Namespace,
	}

	return flag
}

func (o *ApplicationOptions) NewServiceAccountFlag() *cmd.StringFlag {
	flag := &cmd.StringFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "service-account",
			Description: "service account name",
		},
		Option: &o.ServiceAccount,
	}
	return flag
}

func (o *ApplicationOptions) NewUsePreviousFlag() *cmd.BoolFlag {
	flag := &cmd.BoolFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "use-previous",
			Description: fmt.Sprintf("use previous working version (%q)", o.DefaultVersion.PreviousString()),
		},
		Option: &o.UsePrevious,
	}
	return flag
}

func (o *ApplicationOptions) NewVersionFlag() *cmd.StringFlag {
	flag := &cmd.StringFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "version",
			Description: fmt.Sprintf("container image tag (default %q)", o.DefaultVersion.LatestString()),
			Shorthand:   "v",
			Validate: func() error {
				if o.UsePrevious && o.Version != "" {
					return fmt.Errorf("%q flag cannot be used with %q flag", "use-previous", "version")
				}

				if o.UsePrevious {
					o.Version = o.PreviousVersion(*o.Cluster.Version)
					return nil
				}

				if o.Version == "" {
					o.Version = o.LatestVersion(*o.Cluster.Version)
				}

				return nil
			},
		},
		Option: &o.Version,
	}
	return flag
}
