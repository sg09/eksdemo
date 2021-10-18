package kubernetes

import (
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
)

type ResourceManager struct {
	template.Template
}

func (m *ResourceManager) Create(options resource.Options) error {
	manifest, err := m.Render(options)
	if err != nil {
		return err
	}

	return CreateResources(options.GetKubeContext(), manifest)
}

func (m *ResourceManager) Delete(options resource.Options) error {
	return nil
}
