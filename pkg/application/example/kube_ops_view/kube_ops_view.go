package kube_ops_view

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/template"
)

// Codeberg: https://codeberg.org/hjacobs/kube-ops-view
// Manifest: https://codeberg.org/hjacobs/kube-ops-view/src/branch/main/deploy
// Repo:     hjacobs/kube-ops-view

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "kube-ops-view",
			Description: "Kubernetes Operational View",
			Aliases:     []string{"kubeopsview"},
		},

		Installer: &installer.ManifestInstaller{
			AppName: "kube-ops-view",
			ResourceTemplate: &template.TextTemplate{
				Template: deploymentTemplate + rbacTemplate + serviceTemplate,
			},
		},
	}

	app.Options, app.Flags = newOptions()

	return app
}
