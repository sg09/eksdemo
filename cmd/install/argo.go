package install

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/application/argo/argo_cd"

	"github.com/spf13/cobra"
)

var argoApps []func() *application.Application

func NewInstallArgoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "argo",
		Short: "Get stuff done with Kubernetes!",
	}

	// Don't show flag errors for `install argo` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range argoApps {
		cmd.AddCommand(a().NewInstallCmd())
	}

	return cmd
}

func NewUninstallArgoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "argo",
		Short: "Get stuff done with Kubernetes!",
	}

	// Don't show flag errors for `uninstall argo` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range argoApps {
		cmd.AddCommand(a().NewUninstallCmd())
	}

	return cmd
}

func init() {
	argoApps = []func() *application.Application{
		argo_cd.NewApp,
	}
}
