package kustomize

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/kubernetes"
	"eksdemo/pkg/template"
)

type KustomizeInstaller struct {
	ResourceTemplate  template.Template
	KustomizeTemplate template.Template
}

func (i *KustomizeInstaller) Install(options application.Options) error {
	resources, err := i.ResourceTemplate.Render(options)
	if err != nil {
		return err
	}

	kustomization, err := i.KustomizeTemplate.Render(options)
	if err != nil {
		return err
	}

	yaml, err := Kustomize(resources, kustomization)
	if err != nil {
		return err
	}

	err = kubernetes.CreateResources(options.KubeContext(), yaml)
	if err != nil {
		return err
	}

	return nil
}

func (i *KustomizeInstaller) Uninstall(options application.Options) error {
	resources, err := i.ResourceTemplate.Render(options)
	if err != nil {
		return err
	}

	kustomization, err := i.KustomizeTemplate.Render(options)
	if err != nil {
		return err
	}

	yaml, err := Kustomize(resources, kustomization)
	if err != nil {
		return err
	}

	err = kubernetes.DeleteResources(options.KubeContext(), yaml)
	if err != nil {
		return err
	}

	return nil
}
