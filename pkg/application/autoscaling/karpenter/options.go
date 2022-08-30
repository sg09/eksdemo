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
				LatestChart:   "0.16.0",
				Latest:        "v0.16.0",
				PreviousChart: "0.14.0",
				Previous:      "v0.14.0",
			},
		},
	}

	flags = cmd.Flags{}
	return
}
