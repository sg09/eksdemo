package install

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/application/flux/flux_controllers"

	"github.com/spf13/cobra"
)

var fluxApps []func() *application.Application

func NewInstallFluxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flux",
		Short: "GitOps family of projects",
	}

	// Don't show flag errors for `install flux` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range fluxApps {
		cmd.AddCommand(a().NewInstallCmd())
	}

	return cmd
}

func NewUninstallFluxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flux",
		Short: "GitOps family of projects",
	}

	// Don't show flag errors for `uninstall flux` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range fluxApps {
		cmd.AddCommand(a().NewUninstallCmd())
	}

	return cmd
}

func init() {
	fluxApps = []func() *application.Application{
		flux_controllers.NewApp,
	}
}
