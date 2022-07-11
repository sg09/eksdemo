package cilium

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
)

type CiliumOptions struct {
	application.ApplicationOptions

	Wireguard bool
}

func newOptions() (options *CiliumOptions, flags cmd.Flags) {
	options = &CiliumOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace: "kube-system",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.11.6",
				Latest:        "v1.11.6",
				PreviousChart: "1.11.5",
				Previous:      "v1.11.5",
			},
			DisableServiceAccountFlag: true,
		},
		Wireguard: false,
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "wireguard",
				Description: "enable wireguard transparent encryption",
			},
			Option: &options.Wireguard,
		},
	}

	return
}
