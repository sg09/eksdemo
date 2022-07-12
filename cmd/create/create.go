package create

import (
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/acm_certificate"
	"eksdemo/pkg/resource/addon"
	"eksdemo/pkg/resource/amg"
	"eksdemo/pkg/resource/amp"
	"eksdemo/pkg/resource/cluster"
	"eksdemo/pkg/resource/fargate_profile"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/resource/nodegroup"
	"eksdemo/pkg/resource/organization"
	"eksdemo/pkg/resource/servicelb"

	"github.com/spf13/cobra"
)

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create resource(s)",
	}

	// Don't show flag errors for create without a subcommand
	cmd.DisableFlagParsing = true

	cmd.AddCommand(NewAckCmd())
	for _, c := range NewCreateAliasCmds(ack, "ack-") {
		cmd.AddCommand(c)
	}
	cmd.AddCommand(acm_certificate.NewResource().NewCreateCmd())
	cmd.AddCommand(addon.NewResource().NewCreateCmd())
	cmd.AddCommand(amg.NewResource().NewCreateCmd())
	cmd.AddCommand(amp.NewResource().NewCreateCmd())
	cmd.AddCommand(cluster.NewResource().NewCreateCmd())
	cmd.AddCommand(fargate_profile.NewResource().NewCreateCmd())
	cmd.AddCommand(irsa.NewResource().NewCreateCmd())
	cmd.AddCommand(nodegroup.NewResource().NewCreateCmd())
	cmd.AddCommand(nodegroup.NewSpotResource().NewCreateCmd())
	cmd.AddCommand(nodegroup.NewGravitonResource().NewCreateCmd())
	cmd.AddCommand(organization.NewResource().NewCreateCmd())
	cmd.AddCommand(servicelb.NewResource().NewCreateCmd())

	return cmd
}

// This creates alias commands for subcommands under CREATE
func NewCreateAliasCmds(resList []func() *resource.Resource, prefix string) []*cobra.Command {
	cmds := make([]*cobra.Command, 0, len(resList))

	for _, res := range resList {
		r := res()
		r.Command.Name = prefix + r.Command.Name
		r.Command.Hidden = true
		for i, alias := range r.Command.Aliases {
			r.Command.Aliases[i] = prefix + alias
		}
		cmds = append(cmds, r.NewCreateCmd())
	}

	return cmds
}
