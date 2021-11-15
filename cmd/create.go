package cmd

import (
	"eksdemo/pkg/resource/amg"
	"eksdemo/pkg/resource/amp"
	"eksdemo/pkg/resource/cluster"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/resource/nodegroup"
	"eksdemo/pkg/resource/organization"
	"eksdemo/pkg/resource/servicelb"

	"github.com/spf13/cobra"
)

func newCmdCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create resource(s)",
	}

	// Don't show flag errors for create without a subcommand
	cmd.DisableFlagParsing = true

	cmd.AddCommand(amg.NewResource().NewCreateCmd())
	cmd.AddCommand(amp.NewResource().NewCreateCmd())
	cmd.AddCommand(cluster.NewResource().NewCreateCmd())
	cmd.AddCommand(irsa.NewResource().NewCreateCmd())
	cmd.AddCommand(nodegroup.NewResource().NewCreateCmd())
	cmd.AddCommand(nodegroup.NewSpotResource().NewCreateCmd())
	cmd.AddCommand(nodegroup.NewGravitonResource().NewCreateCmd())
	cmd.AddCommand(organization.NewResource().NewCreateCmd())
	cmd.AddCommand(servicelb.NewResource().NewCreateCmd())

	return cmd
}
