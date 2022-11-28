package install

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/application/adot_operator"
	"eksdemo/pkg/application/appmesh_controller"
	"eksdemo/pkg/application/autoscaling/cluster_autoscaler"
	"eksdemo/pkg/application/autoscaling/karpenter"
	"eksdemo/pkg/application/aws_fluentbit"
	"eksdemo/pkg/application/aws_lb_controller"
	"eksdemo/pkg/application/cert_manager"
	"eksdemo/pkg/application/cilium"
	"eksdemo/pkg/application/crossplane"
	"eksdemo/pkg/application/external_dns"
	"eksdemo/pkg/application/falco"
	"eksdemo/pkg/application/grafana_amp"
	"eksdemo/pkg/application/harbor"
	"eksdemo/pkg/application/keycloak_amg"
	"eksdemo/pkg/application/kube_prometheus"
	"eksdemo/pkg/application/metrics_server"
	"eksdemo/pkg/application/prometheus_amp"
	"eksdemo/pkg/application/storage/ebs_csi"
	"eksdemo/pkg/application/velero"

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
	cmd.AddCommand(NewInstallAliasCmds(ack, "ack-")...)
	cmd.AddCommand(adot_operator.NewApp().NewInstallCmd())
	cmd.AddCommand(appmesh_controller.NewApp().NewInstallCmd())
	cmd.AddCommand(NewInstallArgoCmd())
	cmd.AddCommand(NewInstallAliasCmds(argoApps, "argo-")...)
	cmd.AddCommand(NewInstallAutoscalingCmd())
	cmd.AddCommand(NewInstallAliasCmds(autoscalingApps, "autoscaling-")...)
	cmd.AddCommand(NewInstallAliasCmds(autoscalingApps, "as-")...)
	cmd.AddCommand(aws_fluentbit.NewApp().NewInstallCmd())
	cmd.AddCommand(aws_lb_controller.NewApp().NewInstallCmd())
	cmd.AddCommand(cert_manager.NewApp().NewInstallCmd())
	cmd.AddCommand(cilium.NewApp().NewInstallCmd())
	cmd.AddCommand(NewInstallContainerInsightsCmd())
	cmd.AddCommand(NewInstallAliasCmds(containerInsightsApps, "container-insights-")...)
	cmd.AddCommand(NewInstallAliasCmds(containerInsightsApps, "ci-")...)
	cmd.AddCommand(crossplane.NewApp().NewInstallCmd())
	cmd.AddCommand(NewInstallExampleCmd())
	cmd.AddCommand(NewInstallAliasCmds(exampleApps, "example-")...)
	cmd.AddCommand(NewInstallAliasCmds(exampleApps, "ex-")...)
	cmd.AddCommand(external_dns.NewApp().NewInstallCmd())
	cmd.AddCommand(falco.NewApp().NewInstallCmd())
	cmd.AddCommand(NewInstallFluxCmd())
	cmd.AddCommand(NewInstallAliasCmds(fluxApps, "flux-")...)
	cmd.AddCommand(grafana_amp.NewApp().NewInstallCmd())
	cmd.AddCommand(harbor.NewApp().NewInstallCmd())
	cmd.AddCommand(NewInstallIngressCmd())
	cmd.AddCommand(NewInstallAliasCmds(ingressControllers, "ingress-")...)
	cmd.AddCommand(NewInstallIstioCmd())
	cmd.AddCommand(NewInstallAliasCmds(istioApps, "istio-")...)
	cmd.AddCommand(keycloak_amg.NewApp().NewInstallCmd())
	cmd.AddCommand(kube_prometheus.NewApp().NewInstallCmd())
	cmd.AddCommand(NewInstallKubecostCmd())
	cmd.AddCommand(NewInstallAliasCmds(kubecostApps, "kubecost-")...)
	cmd.AddCommand(metrics_server.NewApp().NewInstallCmd())
	cmd.AddCommand(NewInstallPolicyCmd())
	cmd.AddCommand(NewInstallAliasCmds(policyApps, "policy-")...)
	cmd.AddCommand(prometheus_amp.NewApp().NewInstallCmd())
	cmd.AddCommand(NewInstallStorageCmd())
	cmd.AddCommand(NewInstallAliasCmds(storageApps, "storage-")...)
	cmd.AddCommand(velero.NewApp().NewInstallCmd())

	// Hidden commands for popular apps without using the group
	cmd.AddCommand(NewInstallAliasCmds([]func() *application.Application{cluster_autoscaler.NewApp}, "")...)
	cmd.AddCommand(NewInstallAliasCmds([]func() *application.Application{ebs_csi.NewApp}, "")...)
	cmd.AddCommand(NewInstallAliasCmds([]func() *application.Application{karpenter.NewApp}, "")...)

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
