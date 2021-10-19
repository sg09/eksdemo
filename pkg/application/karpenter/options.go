package karpenter

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/util"
	"fmt"
	"log"
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

func (o *KarpenterOptions) PreInstall() error {
	if o.skipSubnetCheck {
		return nil
	}

	if err := util.CheckSubnets(o.ClusterName); err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("\nKarpenter requires subnets tagged with %q to perform subnet discovery\n",
			fmt.Sprintf(util.K8stag, o.ClusterName))
		fmt.Printf("Either run `eksdemo util tag-subnets -c %s` or use the `--skip-subnet-check` flag\n", o.ClusterName)
	}

	log.Fatal("here")

	return nil
}

func (o *KarpenterOptions) PostInstall() error {
	res := karpenterDefaultProvisioner()
	o.AssignCommonResourceOptions(res)
	fmt.Printf("Creating post-install resource: %s\n", res.Common().Name)

	return res.Create()
}
