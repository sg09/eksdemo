package cmd

import (
	"eksdemo/pkg/resource/amg"
	"eksdemo/pkg/resource/amp"
	"eksdemo/pkg/resource/cloudformation"
	"eksdemo/pkg/resource/cluster"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/resource/nodegroup"
	"eksdemo/pkg/resource/organization"

	"github.com/spf13/cobra"
)

func newCmdDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete resource(s)",
	}

	// Don't show flag errors for delete without a subcommand
	cmd.DisableFlagParsing = true

	cmd.AddCommand(amg.NewResource().NewDeleteCmd())
	cmd.AddCommand(amp.NewResource().NewDeleteCmd())
	cmd.AddCommand(cloudformation.NewResource().NewDeleteCmd())
	cmd.AddCommand(cluster.NewResource().NewDeleteCmd())
	cmd.AddCommand(irsa.NewResource().NewDeleteCmd())
	cmd.AddCommand(nodegroup.NewResource().NewDeleteCmd())
	cmd.AddCommand(organization.NewResource().NewDeleteCmd())

	return cmd
}
