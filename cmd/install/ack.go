package install

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/application/ack/ec2_controller"
	"eksdemo/pkg/application/ack/ecr_controller"
	"eksdemo/pkg/application/ack/eks_controller"
	"eksdemo/pkg/application/ack/s3_controller"

	"github.com/spf13/cobra"
)

var ack []func() *application.Application

func NewInstallAckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ack",
		Short: "AWS Controllers for Kubernetes (ACK)",
	}

	// Don't show flag errors for `install ack` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range ack {
		cmd.AddCommand(a().NewInstallCmd())
	}

	return cmd
}

func NewUninstallAckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ack",
		Short: "AWS Controllers for Kubernetes (ACK)",
	}

	// Don't show flag errors for `install ack` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range ack {
		cmd.AddCommand(a().NewUninstallCmd())
	}

	return cmd
}

func init() {
	ack = []func() *application.Application{
		s3_controller.NewApp,
		ec2_controller.NewApp,
		ecr_controller.NewApp,
		eks_controller.NewApp,
	}
}
