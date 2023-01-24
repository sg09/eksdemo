package cmd

import (
	"eksdemo/pkg/resource/acm_certificate"
	"eksdemo/pkg/resource/addon"
	"eksdemo/pkg/resource/amg"
	"eksdemo/pkg/resource/amp_workspace"
	"eksdemo/pkg/resource/cloudformation"
	"eksdemo/pkg/resource/cluster"
	"eksdemo/pkg/resource/dns_record"
	"eksdemo/pkg/resource/ec2_instance"
	"eksdemo/pkg/resource/fargate_profile"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/resource/load_balancer"
	"eksdemo/pkg/resource/log_group"
	"eksdemo/pkg/resource/nodegroup"
	"eksdemo/pkg/resource/organization"
	"eksdemo/pkg/resource/security_group"
	"eksdemo/pkg/resource/target_group"
	"eksdemo/pkg/resource/volume"

	"github.com/spf13/cobra"
)

func newCmdDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete resource(s)",
	}

	// Don't show flag errors for delete without a subcommand
	cmd.DisableFlagParsing = true

	cmd.AddCommand(acm_certificate.NewResource().NewDeleteCmd())
	cmd.AddCommand(addon.NewResource().NewDeleteCmd())
	cmd.AddCommand(amg.NewResource().NewDeleteCmd())
	cmd.AddCommand(amp_workspace.NewResource().NewDeleteCmd())
	cmd.AddCommand(cloudformation.NewResource().NewDeleteCmd())
	cmd.AddCommand(cluster.NewResource().NewDeleteCmd())
	cmd.AddCommand(dns_record.NewResource().NewDeleteCmd())
	cmd.AddCommand(ec2_instance.NewResource().NewDeleteCmd())
	cmd.AddCommand(fargate_profile.NewResource().NewDeleteCmd())
	cmd.AddCommand(irsa.NewResource().NewDeleteCmd())
	cmd.AddCommand(load_balancer.NewResource().NewDeleteCmd())
	cmd.AddCommand(log_group.NewResource().NewDeleteCmd())
	cmd.AddCommand(nodegroup.NewResource().NewDeleteCmd())
	cmd.AddCommand(organization.NewResource().NewDeleteCmd())
	cmd.AddCommand(security_group.NewResource().NewDeleteCmd())
	cmd.AddCommand(target_group.NewResource().NewDeleteCmd())
	cmd.AddCommand(volume.NewResource().NewDeleteCmd())

	return cmd
}
