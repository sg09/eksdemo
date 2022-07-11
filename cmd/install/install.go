package install

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/application/adot_operator"
	"eksdemo/pkg/application/appmesh_controller"
	"eksdemo/pkg/application/aws_fluentbit"
	"eksdemo/pkg/application/aws_lb"
	"eksdemo/pkg/application/cert_manager"
	"eksdemo/pkg/application/cilium"
	"eksdemo/pkg/application/cluster_autoscaler"
	"eksdemo/pkg/application/container_insights"
	"eksdemo/pkg/application/container_insights_prom"
	"eksdemo/pkg/application/ebs_csi"
	"eksdemo/pkg/application/efs_csi"
	"eksdemo/pkg/application/external_dns"
	"eksdemo/pkg/application/falco"
	"eksdemo/pkg/application/fsx_lustre_csi"
	"eksdemo/pkg/application/grafana_amp"
	"eksdemo/pkg/application/karpenter"
	"eksdemo/pkg/application/keycloak_amg"
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
	cmd.AddCommand(aws_fluentbit.NewApp().NewInstallCmd())
	cmd.AddCommand(aws_lb.NewApp().NewInstallCmd())
	cmd.AddCommand(cert_manager.NewApp().NewInstallCmd())
	cmd.AddCommand(cilium.NewApp().NewInstallCmd())
	cmd.AddCommand(cluster_autoscaler.NewApp().NewInstallCmd())
	cmd.AddCommand(container_insights.NewApp().NewInstallCmd())
	cmd.AddCommand(container_insights_prom.NewApp().NewInstallCmd())
	cmd.AddCommand(ebs_csi.NewApp().NewInstallCmd())
	cmd.AddCommand(efs_csi.NewApp().NewInstallCmd())
	cmd.AddCommand(NewInstallExampleCmd())
	for _, c := range NewInstallAliasCmds(exampleApps, "example-") {
		cmd.AddCommand(c)
	}
	cmd.AddCommand(external_dns.NewApp().NewInstallCmd())
	cmd.AddCommand(falco.NewApp().NewInstallCmd())
	cmd.AddCommand(NewInstallFluxCmd())
	for _, c := range NewInstallAliasCmds(fluxApps, "flux-") {
		cmd.AddCommand(c)
	}
	cmd.AddCommand(fsx_lustre_csi.NewApp().NewInstallCmd())
	cmd.AddCommand(grafana_amp.NewApp().NewInstallCmd())
	cmd.AddCommand(NewInstallIngressCmd())
	for _, c := range NewInstallAliasCmds(ingressControllers, "ingress-") {
		cmd.AddCommand(c)
	}
	cmd.AddCommand(NewInstallIstioCmd())
	for _, c := range NewInstallAliasCmds(istioApps, "istio-") {
		cmd.AddCommand(c)
	}
	cmd.AddCommand(karpenter.NewApp().NewInstallCmd())
	cmd.AddCommand(keycloak_amg.NewApp().NewInstallCmd())
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
