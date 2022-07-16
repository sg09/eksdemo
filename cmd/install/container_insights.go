package install

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/application/container_insights/cloudwatch_agent"
	"eksdemo/pkg/application/container_insights/fluentbit"
	"eksdemo/pkg/application/container_insights/prometheus"

	"github.com/spf13/cobra"
)

var containerInsightsApps []func() *application.Application

func NewInstallContainerInsightsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "container-insights",
		Short:   "CloudWatch Container Insights",
		Aliases: []string{"ci"},
	}

	// Don't show flag errors for `install container-insights` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range containerInsightsApps {
		cmd.AddCommand(a().NewInstallCmd())
	}

	return cmd
}

func NewUninstallContainerInsightsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "container-insights",
		Short: "CloudWatch Container Insights",
	}

	// Don't show flag errors for `uninstall container-insights` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range containerInsightsApps {
		cmd.AddCommand(a().NewUninstallCmd())
	}

	return cmd
}

func init() {
	containerInsightsApps = []func() *application.Application{
		cloudwatch_agent.NewApp,
		fluentbit.NewApp,
		prometheus.NewApp,
	}
}
