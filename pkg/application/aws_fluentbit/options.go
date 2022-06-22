package aws_fluentbit

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
)

type FluentbitOptions struct {
	*application.ApplicationOptions
	TailReadFromHead string
}

func addOptions(app *application.Application) *application.Application {
	options := &FluentbitOptions{
		ApplicationOptions: &application.ApplicationOptions{
			Namespace:      "logging",
			ServiceAccount: "fluent-bit",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "0.20.2",
				Latest:        "2.26.0",
				PreviousChart: "0.20.1",
				Previous:      "2.25.1",
			},
		},
		TailReadFromHead: "On",
	}
	app.Options = options

	app.Flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "read-from-head",
				Description: "configuration for tail input plugin",
			},
			Choices: []string{"On", "Off"},
			Option:  &options.TailReadFromHead,
		},
	}
	return app
}
