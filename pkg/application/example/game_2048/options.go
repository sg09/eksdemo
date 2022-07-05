package game_2048

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
)

type Game2048Options struct {
	application.ApplicationOptions

	IngressClass string
	IngressHost  string
	Replicas     int
}

func NewOptions() (options *Game2048Options, flags cmd.Flags) {
	options = &Game2048Options{
		ApplicationOptions: application.ApplicationOptions{
			Namespace: "game-2048",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "latest",
				Previous: "latest",
			},
			DisableServiceAccountFlag: true,
		},
		IngressClass: "alb",
		Replicas:     1,
	}

	flags = cmd.Flags{
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
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "replicas",
				Description: "number of replicas for the deployment",
			},
			Option: &options.Replicas,
		},
	}
	return
}
