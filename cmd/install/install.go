package install

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/application/adot_operator"
	"eksdemo/pkg/application/appmesh_controller"
	"eksdemo/pkg/application/aws_lb"
	"eksdemo/pkg/application/cert_manager"
	"eksdemo/pkg/application/cluster_autoscaler"
	"eksdemo/pkg/application/container_insights"
	"eksdemo/pkg/application/container_insights_prom"
	"eksdemo/pkg/application/ebs_csi"
	"eksdemo/pkg/application/efs_csi"
	"eksdemo/pkg/application/examples/game_2048"
	"eksdemo/pkg/application/external_dns"
	"eksdemo/pkg/application/fluentbit"
	"eksdemo/pkg/application/fsx_lustre_csi"
	"eksdemo/pkg/application/grafana_amp"
	"eksdemo/pkg/application/karpenter"
	"eksdemo/pkg/application/keycloak"
	"eksdemo/pkg/application/kube_prometheus"
	"eksdemo/pkg/application/metrics_server"
	"eksdemo/pkg/application/prometheus_amp"

	"github.com/spf13/cobra"
)

func NewInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install",
		Short:   "install application and any required dependencies",
		Aliases: []string{"inst"},
	}

	// Don't show flag errors for install without a subcommand
	cmd.DisableFlagParsing = true

	cmd.AddCommand(NewInstallAckCmd())
	for _, c := range NewInstallAliasCmds(ack, "ack-") {
		cmd.AddCommand(c)
	}
	cmd.AddCommand(adot_operator.NewApp().NewInstallCmd())
	cmd.AddCommand(appmesh_controller.NewApp().NewInstallCmd())
	cmd.AddCommand(aws_lb.NewApp().NewInstallCmd())
	cmd.AddCommand(cert_manager.NewApp().NewInstallCmd())
	cmd.AddCommand(cluster_autoscaler.NewApp().NewInstallCmd())
	cmd.AddCommand(container_insights.NewApp().NewInstallCmd())
	cmd.AddCommand(container_insights_prom.NewApp().NewInstallCmd())
	cmd.AddCommand(ebs_csi.NewApp().NewInstallCmd())
	cmd.AddCommand(efs_csi.NewApp().NewInstallCmd())
	cmd.AddCommand(external_dns.NewApp().NewInstallCmd())
	cmd.AddCommand(fluentbit.NewApp().NewInstallCmd())
	cmd.AddCommand(game_2048.NewApp().NewInstallCmd())
	cmd.AddCommand(grafana_amp.NewApp().NewInstallCmd())
	cmd.AddCommand(fsx_lustre_csi.NewApp().NewInstallCmd())
	cmd.AddCommand(NewInstallIstioCmd())
	for _, c := range NewInstallAliasCmds(istioCmds, "istio-") {
		cmd.AddCommand(c)
	}
	cmd.AddCommand(karpenter.NewApp().NewInstallCmd())
	cmd.AddCommand(keycloak.NewApp().NewInstallCmd())
	cmd.AddCommand(kube_prometheus.NewApp().NewInstallCmd())
	cmd.AddCommand(metrics_server.NewApp().NewInstallCmd())
	cmd.AddCommand(prometheus_amp.NewApp().NewInstallCmd())

	return cmd
}

// This creates alias commands for subcommands under INSTALL
func NewInstallAliasCmds(appList []func() *application.Application, prefix string) []*cobra.Command {
	cmds := make([]*cobra.Command, 0, len(appList))

	for _, app := range appList {
		a := app()
		a.Command.Name = prefix + a.Command.Name
		a.Command.Hidden = true
		for i, alias := range a.Command.Aliases {
			a.Command.Aliases[i] = prefix + alias
		}
		cmds = append(cmds, a.NewInstallCmd())
	}

	return cmds
}
