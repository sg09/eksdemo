package argo_cd

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
)

type ArgoCdOptions struct {
	application.ApplicationOptions

	AdminPassword string
	IngressClass  string
	IngressHost   string
}

func newOptions() (options *ArgoCdOptions, flags cmd.Flags) {
	options = &ArgoCdOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "4.9.14",
				Latest:        "v2.4.6",
				PreviousChart: "4.9.12",
				Previous:      "v2.4.4",
			},
			DisableServiceAccountFlag: true,
			Namespace:                 "argocd",
		},
		IngressClass: "alb",
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "admin-pass",
				Description: "Argo CD admin password",
				Required:    true,
				Shorthand:   "P",
			},
			Option: &options.AdminPassword,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ingress-class",
				Description: "name of IngressClass",
			},
			Option: &options.IngressClass,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ingress-host",
				Description: "hostname for Ingress with TLS (requires Ingress Controller and ExternalDNS)",
				Shorthand:   "I",
			},
			Option: &options.IngressHost,
		},
	}

	return
}
