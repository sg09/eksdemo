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
				Latest:   "v0.8.2",
				Previous: "v0.8.2",
			},
		},
	}

	flags = cmd.Flags{}
	return
}
