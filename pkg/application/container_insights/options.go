package container_insights

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/template"
)

type ContainerInsightsOptions struct {
	*application.ApplicationOptions
	FluentBit      bool
	HttpServer     string
	HttpServerPort string
	ReadFromHead   string
	ReadFromTail   string
}

func addOptions(app *application.Application) *application.Application {
	options := &ContainerInsightsOptions{
		FluentBit:      false,
		HttpServer:     "On",
		HttpServerPort: "2020",
		ReadFromHead:   "Off",
		ReadFromTail:   "On",
		ApplicationOptions: &application.ApplicationOptions{
			DisableNamespaceFlag:      true,
			DisableServiceAccountFlag: true,
			Namespace:                 "amazon-cloudwatch",
			ServiceAccount:            "cloudwatch-agent",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "1.247348.0b251302",
				Previous: "1.247348.0b251302",
			},
		},
	}
	app.Options = options

	app.Flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "include-logging",
				Description: "include vended Fluent Bit configuration",
			},
			Option: &options.FluentBit,
		},
	}
	return app
}

func (o *ContainerInsightsOptions) PostInstall(_ string, _ []*resource.Resource) error { //include app in params?
	if !o.FluentBit {
		return nil
	}

	o.ServiceAccount = "fluent-bit"

	fluentBit := &application.Application{
		Options: o,

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "fluent-bit-irsa",
				},
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"},
			}),
		},

		Installer: &installer.KustomizeInstaller{
			ResourceTemplate: &template.TextTemplate{
				Template: fluentBitManifestTemplate,
			},
			KustomizeTemplate: &template.TextTemplate{
				Template: fluentBitKustomizeTemplate,
			},
		},
	}

	if err := fluentBit.CreateDependencies(); err != nil {
		return err
	}

	return fluentBit.Install(fluentBit.Options)
}
