package karpenter

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"strings"

	"github.com/spf13/cobra"
)

type KarpenterOptions struct {
	application.ApplicationOptions

	AMIFamily            string
	Replicas             int
	TTLSecondsAfterEmpty int
}

func newOptions() (options *KarpenterOptions, flags cmd.Flags) {
	options = &KarpenterOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "karpenter",
			ServiceAccount: "karpenter",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "v0.22.1",
				Latest:        "v0.22.1",
				PreviousChart: "v0.21.1",
				Previous:      "v0.21.1",
			},
		},
		AMIFamily: "AL2",
		Replicas:  1,
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ami-family",
				Description: "provisioner ami family",
				Shorthand:   "A",
				Validate: func(cmd *cobra.Command, args []string) error {
					if strings.EqualFold(options.AMIFamily, "Al2") {
						options.AMIFamily = "AL2"
						return nil
					}
					if strings.EqualFold(options.AMIFamily, "Bottlerocket") {
						options.AMIFamily = "Bottlerocket"
						return nil
					}
					if strings.EqualFold(options.AMIFamily, "Ubuntu") {
						options.AMIFamily = "Ubuntu"
						return nil
					}
					return nil
				},
			},
			Option:  &options.AMIFamily,
			Choices: []string{"AL2", "Bottlerocket", "Ubuntu"},
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "replicas",
				Description: "number of replicas for the controller deployment",
			},
			Option: &options.Replicas,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ttl-after-empty",
				Description: "provisioner ttl seconds after empty (disables consolidation)",
				Shorthand:   "T",
			},
			Option: &options.TTLSecondsAfterEmpty,
		},
	}
	return
}
