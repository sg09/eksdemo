package kubernetes

import (
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
	"fmt"
)

type ResourceManager struct {
	Resource string
	template.Template
}

func (m *ResourceManager) Create(options resource.Options) error {
	manifest, err := m.Render(options)
	if err != nil {
		return err
	}

	fmt.Println(manifest)

	// return CreateResources(options.KubeContext(), yaml)
	return nil
}

func (m *ResourceManager) Delete(options resource.Options) error {
	return nil
}
