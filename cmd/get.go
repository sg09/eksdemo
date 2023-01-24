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
	"eksdemo/pkg/resource/amp_workspace"
	"eksdemo/pkg/resource/auto_scaling_group"
	"eksdemo/pkg/resource/availability_zone"
	"eksdemo/pkg/resource/cloudformation"
	"eksdemo/pkg/resource/cloudtrail_event"
	"eksdemo/pkg/resource/cloudtrail_trail"
	"eksdemo/pkg/resource/cluster"
	"eksdemo/pkg/resource/dns_record"
	"eksdemo/pkg/resource/ec2_instance"
	"eksdemo/pkg/resource/ecr_repository"
	"eksdemo/pkg/resource/elastic_ip"
	"eksdemo/pkg/resource/event_rule"
	"eksdemo/pkg/resource/fargate_profile"
	"eksdemo/pkg/resource/hosted_zone"
	"eksdemo/pkg/resource/iam_oidc"
	"eksdemo/pkg/resource/iam_policy"
	"eksdemo/pkg/resource/iam_role"
	"eksdemo/pkg/resource/internet_gateway"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/resource/kms_key"
	"eksdemo/pkg/resource/listener"
	"eksdemo/pkg/resource/listener_rule"
	"eksdemo/pkg/resource/load_balancer"
	"eksdemo/pkg/resource/log_event"
	"eksdemo/pkg/resource/log_group"
	"eksdemo/pkg/resource/log_stream"
	"eksdemo/pkg/resource/metric"
	"eksdemo/pkg/resource/nat_gateway"
	"eksdemo/pkg/resource/network_acl"
	"eksdemo/pkg/resource/network_acl_rule"
	"eksdemo/pkg/resource/network_interface"
	"eksdemo/pkg/resource/node"
	"eksdemo/pkg/resource/nodegroup"
	"eksdemo/pkg/resource/organization"
	"eksdemo/pkg/resource/route_table"
	"eksdemo/pkg/resource/s3_bucket"
	"eksdemo/pkg/resource/security_group"
	"eksdemo/pkg/resource/security_group_rule"
	"eksdemo/pkg/resource/sqs_queue"
	"eksdemo/pkg/resource/ssm_node"
	"eksdemo/pkg/resource/ssm_session"
	"eksdemo/pkg/resource/subnet"
	"eksdemo/pkg/resource/target_group"
	"eksdemo/pkg/resource/target_health"
	"eksdemo/pkg/resource/volume"
	"eksdemo/pkg/resource/vpc"
	"eksdemo/pkg/resource/vpc_endpoint"
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
		cluster, err := aws.NewEKSClient().DescribeCluster(clusterName)
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
	getCmd.AddCommand(amp_workspace.NewResource().NewGetCmd())
	getCmd.AddCommand(auto_scaling_group.NewResource().NewGetCmd())
	getCmd.AddCommand(availability_zone.NewResource().NewGetCmd())
	getCmd.AddCommand(cloudformation.NewResource().NewGetCmd())
	getCmd.AddCommand(cloudtrail_event.NewResource().NewGetCmd())
	getCmd.AddCommand(cloudtrail_trail.NewResource().NewGetCmd())
	getCmd.AddCommand(cluster.NewResource().NewGetCmd())
	getCmd.AddCommand(dns_record.NewResource().NewGetCmd())
	getCmd.AddCommand(ec2_instance.NewResource().NewGetCmd())
	getCmd.AddCommand(ecr_repository.NewResource().NewGetCmd())
	getCmd.AddCommand(elastic_ip.NewResource().NewGetCmd())
	getCmd.AddCommand(event_rule.NewResource().NewGetCmd())
	getCmd.AddCommand(fargate_profile.NewResource().NewGetCmd())
	getCmd.AddCommand(hosted_zone.NewResource().NewGetCmd())
	getCmd.AddCommand(iam_oidc.NewResource().NewGetCmd())
	getCmd.AddCommand(iam_policy.NewResource().NewGetCmd())
	getCmd.AddCommand(iam_role.NewResource().NewGetCmd())
	getCmd.AddCommand(internet_gateway.NewResource().NewGetCmd())
	getCmd.AddCommand(irsa.NewResource().NewGetCmd())
	getCmd.AddCommand(kms_key.NewResource().NewGetCmd())
	getCmd.AddCommand(listener.NewResource().NewGetCmd())
	getCmd.AddCommand(listener_rule.NewResource().NewGetCmd())
	getCmd.AddCommand(load_balancer.NewResource().NewGetCmd())
	getCmd.AddCommand(log_event.NewResource().NewGetCmd())
	getCmd.AddCommand(log_group.NewResource().NewGetCmd())
	getCmd.AddCommand(log_stream.NewResource().NewGetCmd())
	getCmd.AddCommand(metric.NewResource().NewGetCmd())
	getCmd.AddCommand(nat_gateway.NewResource().NewGetCmd())
	getCmd.AddCommand(network_acl.NewResource().NewGetCmd())
	getCmd.AddCommand(network_acl_rule.NewResource().NewGetCmd())
	getCmd.AddCommand(network_interface.NewResource().NewGetCmd())
	getCmd.AddCommand(node.NewResource().NewGetCmd())
	getCmd.AddCommand(nodegroup.NewResource().NewGetCmd())
	getCmd.AddCommand(organization.NewResource().NewGetCmd())
	getCmd.AddCommand(route_table.NewResource().NewGetCmd())
	getCmd.AddCommand(s3_bucket.NewResource().NewGetCmd())
	getCmd.AddCommand(security_group.NewResource().NewGetCmd())
	getCmd.AddCommand(security_group_rule.NewResource().NewGetCmd())
	getCmd.AddCommand(sqs_queue.NewResource().NewGetCmd())
	getCmd.AddCommand(ssm_node.NewResource().NewGetCmd())
	getCmd.AddCommand(ssm_session.NewResource().NewGetCmd())
	getCmd.AddCommand(subnet.NewResource().NewGetCmd())
	getCmd.AddCommand(target_group.NewResource().NewGetCmd())
	getCmd.AddCommand(target_health.NewResource().NewGetCmd())
	getCmd.AddCommand(volume.NewResource().NewGetCmd())
	getCmd.AddCommand(vpc.NewResource().NewGetCmd())
	getCmd.AddCommand(vpc_endpoint.NewResource().NewGetCmd())

	// Don't show flag errors for install without a subcommand
	getCmd.DisableFlagParsing = true

	getHelmCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "cluster (required)")
	getHelmCmd.MarkFlagRequired("cluster")
	getHelmCmd.Flags().VarP(cmd.NewOutputFlag(&output), "output", "o", "output format (json|table|yaml)")
	getCmd.AddCommand(getHelmCmd)
}
