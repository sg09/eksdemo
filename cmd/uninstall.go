package cmd

import (
	"eksdemo/pkg/application/aws_lb"
	"eksdemo/pkg/application/cluster_autoscaler"
	"eksdemo/pkg/application/container_insights"
	"eksdemo/pkg/application/container_insights_prom"
	"eksdemo/pkg/application/ebs_csi"
	"eksdemo/pkg/application/efs_csi"
	"eksdemo/pkg/application/external_dns"
	"eksdemo/pkg/application/fluentbit"
	"eksdemo/pkg/application/fsx_lustre_csi"
	"eksdemo/pkg/application/karpenter"
	"eksdemo/pkg/application/kube_prometheus"
	"eksdemo/pkg/application/metrics_server"
	"eksdemo/pkg/application/prometheus_amp"

	"github.com/spf13/cobra"
)

func newCmdUninstall() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "uninstall",
		Short:   "uninstall application and delete dependencies",
		Aliases: []string{"uninst"},
	}

	// Don't show flag errors for uninstall without a subcommand
	cmd.DisableFlagParsing = true

	cmd.AddCommand(aws_lb.NewApp().NewUninstallCmd())
	cmd.AddCommand(cluster_autoscaler.NewApp().NewUninstallCmd())
	cmd.AddCommand(container_insights.NewApp().NewUninstallCmd())
	cmd.AddCommand(container_insights_prom.NewApp().NewUninstallCmd())
	cmd.AddCommand(ebs_csi.NewApp().NewUninstallCmd())
	cmd.AddCommand(efs_csi.NewApp().NewUninstallCmd())
	cmd.AddCommand(external_dns.NewApp().NewUninstallCmd())
	cmd.AddCommand(fluentbit.NewApp().NewUninstallCmd())
	cmd.AddCommand(fsx_lustre_csi.NewApp().NewUninstallCmd())
	cmd.AddCommand(karpenter.NewApp().NewUninstallCmd())
	cmd.AddCommand(kube_prometheus.NewApp().NewUninstallCmd())
	cmd.AddCommand(metrics_server.NewApp().NewUninstallCmd())
	cmd.AddCommand(prometheus_amp.NewApp().NewUninstallCmd())

	return cmd
}
