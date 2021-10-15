package fluentbit

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"fmt"
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
				Latest:   "2.19.1",
				Previous: "2.18.0",
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
				Validate: func() error {
					if options.TailReadFromHead == "On" || options.TailReadFromHead == "Off" {
						return nil
					}
					return fmt.Errorf("flag %q must be either %q or %q", "read-from-head", "On", "Off")
				},
			},
			Choices: []string{"On", "Off"},
			Option:  &options.TailReadFromHead,
		},
	}
	return app
}
