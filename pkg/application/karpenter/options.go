package karpenter

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/util"
	"errors"
	"fmt"
)

type KarpenterOptions struct {
	application.ApplicationOptions

	skipSubnetCheck bool
}

func NewOptions() (options *KarpenterOptions, flags cmd.Flags) {
	options = &KarpenterOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "karpenter",
			ServiceAccount: "karpenter",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "v0.4.0",
				Previous: "v0.3.4",
			},
		},
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "skip-subnet-check",
				Description: "don't check subnets for required Karpenter tags",
			},
			Option: &options.skipSubnetCheck,
		},
	}
	return
}

func (o *KarpenterOptions) PreDependencies() error {
	if o.skipSubnetCheck {
		return nil
	}

	if err := util.CheckSubnets(o.ClusterName); err != nil {
		errMsg := err.Error()
		errMsg += fmt.Sprintf("\n\nKarpenter requires subnets tagged with %q to perform subnet discovery\n",
			fmt.Sprintf(util.K8stag, o.ClusterName))
		errMsg += fmt.Sprintf("Either run `eksdemo util tag-subnets -c %s` or use the `--skip-subnet-check` flag", o.ClusterName)
		return errors.New(errMsg)
	}

	return nil
}

func (o *KarpenterOptions) PostInstall() error {
	res := karpenterDefaultProvisioner()
	o.AssignCommonResourceOptions(res)
	fmt.Printf("Creating post-install resource: %s\n", res.Common().Name)

	return res.Create()
}
