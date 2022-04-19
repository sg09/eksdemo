package install

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/application/example/eks_workshop"
	"eksdemo/pkg/application/example/game_2048"
	"eksdemo/pkg/application/example/inflate"

	"github.com/spf13/cobra"
)

var exampleApps []func() *application.Application

func NewInstallExampleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "example",
		Short: "Example Applications",
	}

	// Don't show flag errors for `install example` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range exampleApps {
		cmd.AddCommand(a().NewInstallCmd())
	}

	return cmd
}

func NewUninstallExampleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "example",
		Short: "Example Applications",
	}

	// Don't show flag errors for `install example` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range exampleApps {
		cmd.AddCommand(a().NewUninstallCmd())
	}

	return cmd
}

func init() {
	exampleApps = []func() *application.Application{
		eks_workshop.NewApp,
		game_2048.NewApp,
		inflate.NewApp,
	}
}
