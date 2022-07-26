package install

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/application/autoscaling/keda"

	"github.com/spf13/cobra"
)

var autoscalingApps []func() *application.Application

func NewInstallAutoscalingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "autoscaling",
		Short: "Kubernetes Autoscaling Applications",
	}

	// Don't show flag errors for `install autoscaling` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range autoscalingApps {
		cmd.AddCommand(a().NewInstallCmd())
	}

	return cmd
}

func NewUninstallAutoscalingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "autoscaling",
		Short: "Kubernetes Autoscaling Applications",
	}

	// Don't show flag errors for `uninstall autoscaling` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range autoscalingApps {
		cmd.AddCommand(a().NewUninstallCmd())
	}

	return cmd
}

func init() {
	autoscalingApps = []func() *application.Application{
		keda.NewApp,
	}
}
