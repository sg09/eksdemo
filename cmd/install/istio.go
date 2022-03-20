package install

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/application/istio/istio_base"
	"eksdemo/pkg/application/istio/istiod"

	"github.com/spf13/cobra"
)

var istioCmds []func() *application.Application

func NewInstallIstioCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "istio",
		Short: "Istio Service Mesh",
	}

	// Don't show flag errors for `install ack` without a subcommand
	cmd.DisableFlagParsing = true

	for _, i := range istioCmds {
		cmd.AddCommand(i().NewInstallCmd())
	}

	return cmd
}

func NewUninstallIstioCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "istio",
		Short: "Istio Service Mesh",
	}

	// Don't show flag errors for `install istio` without a subcommand
	cmd.DisableFlagParsing = true

	for _, i := range istioCmds {
		cmd.AddCommand(i().NewUninstallCmd())
	}

	return cmd
}

func init() {
	istioCmds = []func() *application.Application{
		// bookinfo.NewApp,
		istio_base.NewApp,
		// istio_egress.NewApp,
		// istio_ingress.NewApp,
		istiod.NewApp,
		// kiali.NewApp,
	}
}
