package cmd

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/helm"
	"eksdemo/pkg/kubernetes"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource/acm_certificate"
	"eksdemo/pkg/resource/addon"
	"eksdemo/pkg/resource/amg"
	"eksdemo/pkg/resource/amp"
	"eksdemo/pkg/resource/auto_scaling_group"
	"eksdemo/pkg/resource/cloudformation"
	"eksdemo/pkg/resource/cluster"
	"eksdemo/pkg/resource/dns_record"
	"eksdemo/pkg/resource/ec2_instance"
	"eksdemo/pkg/resource/fargate_profile"
	"eksdemo/pkg/resource/hosted_zone"
	"eksdemo/pkg/resource/iam_role"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/resource/load_balancer"
	"eksdemo/pkg/resource/log_group"
	"eksdemo/pkg/resource/network_interface"
	"eksdemo/pkg/resource/node"
	"eksdemo/pkg/resource/nodegroup"
	"eksdemo/pkg/resource/organization"
	"eksdemo/pkg/resource/s3_bucket"
	"eksdemo/pkg/resource/security_group"
	"eksdemo/pkg/resource/security_group_rule"
	"eksdemo/pkg/resource/subnet"
	"eksdemo/pkg/resource/target_group"
	"eksdemo/pkg/resource/target_health"
	"eksdemo/pkg/resource/volume"
	"eksdemo/pkg/resource/vpc"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var output printer.Output
var clusterName string

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "view resource(s)",
}

var getHelmCmd = &cobra.Command{
	Use:     "application",
	Short:   "Installed Applications",
	Aliases: []string{"app", "apps", "helm"},
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster, err := aws.EksDescribeCluster(clusterName)
		if err != nil {
			return err
		}

		kubeContext, err := kubernetes.KubeContextForCluster(cluster)
		if err != nil {
			return err
		}
		if kubeContext == "" {
			return fmt.Errorf("cluster \"%s\" not found in Kubeconfig", clusterName)
		}

		releases, err := helm.List(kubeContext)
		if err != nil {
			return err
		}

		return output.Print(os.Stdout, printer.NewApplicationPrinter(releases))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.AddCommand(acm_certificate.NewResource().NewGetCmd())
	getCmd.AddCommand(addon.NewResource().NewGetCmd())
	getCmd.AddCommand(addon.NewVersionsResource().NewGetCmd())
	getCmd.AddCommand(amg.NewResource().NewGetCmd())
	getCmd.AddCommand(amp.NewResource().NewGetCmd())
	getCmd.AddCommand(auto_scaling_group.NewResource().NewGetCmd())
	getCmd.AddCommand(cloudformation.NewResource().NewGetCmd())
	getCmd.AddCommand(cluster.NewResource().NewGetCmd())
	getCmd.AddCommand(dns_record.NewResource().NewGetCmd())
	getCmd.AddCommand(ec2_instance.NewResource().NewGetCmd())
	getCmd.AddCommand(fargate_profile.NewResource().NewGetCmd())
	getCmd.AddCommand(hosted_zone.NewResource().NewGetCmd())
	getCmd.AddCommand(iam_role.NewResource().NewGetCmd())
	getCmd.AddCommand(irsa.NewResource().NewGetCmd())
	getCmd.AddCommand(load_balancer.NewResource().NewGetCmd())
	getCmd.AddCommand(log_group.NewResource().NewGetCmd())
	getCmd.AddCommand(network_interface.NewResource().NewGetCmd())
	getCmd.AddCommand(node.NewResource().NewGetCmd())
	getCmd.AddCommand(nodegroup.NewResource().NewGetCmd())
	getCmd.AddCommand(organization.NewResource().NewGetCmd())
	getCmd.AddCommand(s3_bucket.NewResource().NewGetCmd())
	getCmd.AddCommand(security_group.NewResource().NewGetCmd())
	getCmd.AddCommand(security_group_rule.NewResource().NewGetCmd())
	getCmd.AddCommand(subnet.NewResource().NewGetCmd())
	getCmd.AddCommand(target_group.NewResource().NewGetCmd())
	getCmd.AddCommand(target_health.NewResource().NewGetCmd())
	getCmd.AddCommand(volume.NewResource().NewGetCmd())
	getCmd.AddCommand(vpc.NewResource().NewGetCmd())

	// Don't show flag errors for install without a subcommand
	getCmd.DisableFlagParsing = true

	getHelmCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "cluster (required)")
	getHelmCmd.MarkFlagRequired("cluster")
	getHelmCmd.Flags().VarP(cmd.NewOutputFlag(&output), "output", "o", "output format (json|table|yaml)")
	getCmd.AddCommand(getHelmCmd)
}
