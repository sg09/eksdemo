package karpenter

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
)

type KarpenterOptions struct {
	application.ApplicationOptions
}

func NewOptions() (options *KarpenterOptions, flags cmd.Flags) {
	options = &KarpenterOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "karpenter",
			ServiceAccount: "karpenter",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "0.14.0",
				Latest:        "v0.14.0",
				PreviousChart: "0.13.2",
				Previous:      "v0.13.2",
			},
		},
	}

	flags = cmd.Flags{}
	return
}
