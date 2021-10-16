package eksctl

import (
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
)

const EksctlHeader = `
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: {{ .ClusterName }}
  region: {{ .Region }}
  tags: 
    eksdemo: v0.1.0-alpha
{{- if .KubernetesVersion }}
  version: {{ .KubernetesVersion | printf "%q" }}
{{ end }}
`

type ResourceManager struct {
	Resource      string
	Template      template.Template
	ApproveCreate bool
	ApproveDelete bool
}

func (e *ResourceManager) Create(options resource.Options) error {
	eksctlConfig, err := e.Template.Render(options)

	if err != nil {
		return err
	}

	args := []string{
		"create",
		e.Resource,
		"-f",
		"-",
	}

	if e.ApproveCreate {
		args = append(args, "--approve")
	}

	return Command(args, eksctlConfig)
}

func (e *ResourceManager) Delete(options resource.Options) error {
	options.PrepForDelete()
	eksctlConfig, err := e.Template.Render(options)

	if err != nil {
		return err
	}

	args := []string{
		"delete",
		e.Resource,
		"-f",
		"-",
	}

	if e.ApproveDelete {
		args = append(args, "--approve")
	}

	return Command(args, eksctlConfig)
}
