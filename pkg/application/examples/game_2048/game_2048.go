package game_2048

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/template"
)

// GitHub:   https://github.com/alexwhen/docker-2048
// Manifest: https://github.com/kubernetes-sigs/aws-load-balancer-controller/blob/main/docs/examples/2048/2048_full_latest.yaml
// Repo:     https://gallery.ecr.aws/l6m2t8p7/docker-2048

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "game-2048",
			Description: "Example Game 2048",
			Aliases:     []string{"game2048", "2048"},
		},

		Installer: &installer.ManifestInstaller{
			AppName: "game-2048",
			ResourceTemplate: &template.TextTemplate{
				Template: gameManifestTemplate,
			},
		},

		Options: &application.ApplicationOptions{
			Namespace: "game-2048",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "latest",
				Previous: "latest",
			},
		},
	}

	app.Options, app.Flags = NewOptions()

	return app
}
