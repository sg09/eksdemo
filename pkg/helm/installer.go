package helm

import (
	"bytes"
	"context"
	"eksdemo/pkg/application"
	"eksdemo/pkg/kubernetes"
	"eksdemo/pkg/kustomize"
	"eksdemo/pkg/template"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

type HelmInstaller struct {
	ChartName           string
	DryRun              bool
	PostRenderKustomize template.Template
	PVCLabels           map[string]string
	ReleaseName         string
	RepositoryURL       string
	ValuesTemplate      template.Template
	VersionField        string
	Wait                bool
	application.Options
}

func (h *HelmInstaller) Install(options application.Options) error {

	valuesFile, err := h.ValuesTemplate.Render(options)
	if err != nil {
		return err
	}

	ic := &InstallConfiguration{
		AppVersion:    options.Common().Version,
		ChartName:     h.ChartName,
		Namespace:     options.Common().Namespace,
		ReleaseName:   h.ReleaseName,
		RepositoryURL: h.RepositoryURL,
		ValuesFile:    valuesFile,
		Wait:          h.Wait,
	}

	if h.DryRun {
		fmt.Println("\nHelm Installer Dry Run:")
		fmt.Printf("%+v\n", ic)

		if h.PostRenderKustomize != nil {
			kustomization, err := h.PostRenderKustomize.Render(options)
			if err != nil {
				return err
			}
			fmt.Println("\nHelm Installer Post Render Kustomize Dry Run:")
			fmt.Printf("%s\n", kustomization)
		}
		return nil
	}

	if h.PostRenderKustomize != nil {
		h.Options = options
		ic.PostRenderer = h
	}

	return Install(ic, options.KubeContext())
}

func (h *HelmInstaller) SetDryRun() {
	h.DryRun = true
}

func (h *HelmInstaller) Uninstall(options application.Options) error {
	o := options.Common()

	fmt.Printf("Checking status of Helm release: %s, in namespace: %s\n", h.ReleaseName, o.Namespace)
	if _, err := Status(o.KubeContext(), h.ReleaseName, o.Namespace); err != nil {
		return err
	}

	fmt.Println("Status validated. Uninstalling...")
	err := Uninstall(o.KubeContext(), h.ReleaseName, o.Namespace)
	if err != nil {
		return err
	}

	if len(h.PVCLabels) == 0 {
		return nil
	}

	// Delete any leftover PVCs as `helm uninstall` won't delete them
	// https://github.com/helm/helm/issues/5156
	client, err := kubernetes.Client(o.KubeContext())
	if err != nil {
		return fmt.Errorf("failed creating kubernetes client: %w", err)
	}

	selector := labels.NewSelector()

	for k, v := range h.PVCLabels {
		req, err := labels.NewRequirement(k, selection.Equals, []string{v})
		if err != nil {
			return err
		}
		selector = selector.Add(*req)
	}

	fmt.Printf("Deleting PVCs with labels: %s\n", selector.String())

	return client.CoreV1().PersistentVolumeClaims(o.Namespace).DeleteCollection(context.Background(),
		metav1.DeleteOptions{},
		metav1.ListOptions{
			LabelSelector: selector.String(),
		},
	)
}

// PostRender
func (h *HelmInstaller) Run(renderedManifests *bytes.Buffer) (modifiedManifests *bytes.Buffer, err error) {
	kustomization, err := h.PostRenderKustomize.Render(h.Options)
	if err != nil {
		return nil, err
	}

	yaml, err := kustomize.Kustomize(renderedManifests.String(), kustomization)
	if err != nil {
		return nil, err
	}

	return bytes.NewBufferString(yaml), nil
}
