package cilium

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/installer"
	"eksdemo/pkg/template"
)

// Docs:    https://docs.cilium.io/
// GitHub:  https://github.com/cilium/cilium
// Helm:    https://github.com/cilium/cilium/tree/master/install/kubernetes/cilium
// Repo:    https://quay.io/repository/cilium/cilium
// Version: Latest is v1.11.6 (as of 07/10/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "cilium",
			Description: "eBPF-based Networking, Observability, Security",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "cilium",
			ReleaseName:   "cilium",
			RepositoryURL: "https://helm.cilium.io/",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options, app.Flags = newOptions()

	return app
}

const valuesTemplate = `---
image:
  tag: {{ .Version }}
cni:
  # -- Install the CNI configuration and binary files into the filesystem.
  install: true
  # -- Configure chaining on top of other CNI plugins. Possible values:
  chainingMode: aws-cni
{{- if .Wireguard }}
encryption:
  # -- Enable transparent network encryption.
  enabled: true
  # -- Encryption method. Can be either ipsec or wireguard.
  type: wireguard
{{- end }}
enableIPv4Masquerade: false
tunnel: disabled
{{- if .Wireguard }}
l7Proxy: false
{{- end }}
`
