package install

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/application/autoscaling/cluster_autoscaler"
	"eksdemo/pkg/application/autoscaling/goldilocks"
	"eksdemo/pkg/application/autoscaling/inflate"
	"eksdemo/pkg/application/autoscaling/karpenter"
	"eksdemo/pkg/application/autoscaling/keda"
	"eksdemo/pkg/application/autoscaling/vpa"

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
		cluster_autoscaler.NewApp,
		goldilocks.NewApp,
		inflate.NewApp,
		karpenter.NewApp,
		keda.NewApp,
		vpa.NewApp,
	}
}
