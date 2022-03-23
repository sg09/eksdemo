package game_2048

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
)

type Game2048Options struct {
	application.ApplicationOptions

	Replicas int
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
		Replicas: 1,
	}

	flags = cmd.Flags{
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
