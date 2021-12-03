package create

import (
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/ack/s3"

	"github.com/spf13/cobra"
)

var ack []func() *resource.Resource

func NewAckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ack",
		Short: "AWS Controllers for Kubernetes (ACK)",
	}

	// Don't show flag errors for `create ack`` without a subcommand
	cmd.DisableFlagParsing = true

	for _, r := range ack {
		cmd.AddCommand(r().NewCreateCmd())
	}

	return cmd
}

func init() {
	ack = []func() *resource.Resource{
		s3.NewResource,
	}
}
