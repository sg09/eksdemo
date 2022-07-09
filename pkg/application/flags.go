package application

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/kubernetes"
	"fmt"
)

func (o *ApplicationOptions) NewChartVersionFlag() *cmd.StringFlag {
	flag := &cmd.StringFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "chart-version",
			Description: fmt.Sprintf("chart version (default %q)", o.DefaultVersion.LatestChartVersion()),
			Validate: func() error {
				if o.UsePrevious && o.ChartVersion != "" {
					return fmt.Errorf("%q flag cannot be used with %q flag", "use-previous", "chart-version")
				}

				if o.UsePrevious {
					o.ChartVersion = o.PreviousChartVersion()
					return nil
				}

				if o.ChartVersion == "" {
					o.ChartVersion = o.LatestChartVersion()
				}

				return nil
			},
		},
		Option: &o.ChartVersion,
	}
	return flag
}

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

func (o *ApplicationOptions) NewSetFlag() *cmd.StringSliceFlag {
	flag := &cmd.StringSliceFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "set",
			Description: "set chart values (can specify multiple or separate values with commas: key1=val1,key2=val2)",
		},
		Option: &o.SetValues,
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
			Name: "use-previous",
			Description: fmt.Sprintf("use previous working chart/app versions (%q/%q)",
				o.DefaultVersion.PreviousChartVersion(), o.DefaultVersion.PreviousString()),
		},
		Option: &o.UsePrevious,
	}
	return flag
}

func (o *ApplicationOptions) NewVersionFlag() *cmd.StringFlag {
	flag := &cmd.StringFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "version",
			Description: fmt.Sprintf("application version (default %q)", o.DefaultVersion.LatestString()),
			Shorthand:   "v",
			Validate: func() error {
				if o.UsePrevious && o.Version != "" {
					return fmt.Errorf("%q flag cannot be used with %q flag", "use-previous", "version")
				}

				if o.LockVersionFlag {
					return fmt.Errorf("version is locked and cannot be changed")
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
