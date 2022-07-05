package kubernetes

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func CreateResources(kubeContext, manifest string) error {
	getter := genericclioptions.NewConfigFlags(true)
	getter.Context = &kubeContext
	factory := cmdutil.NewFactory(getter)
	builder := factory.NewBuilder().Unstructured()

	switch {
	case strings.Index(manifest, "http://") == 0 || strings.Index(manifest, "https://") == 0:
		url, err := url.Parse(manifest)
		if err != nil {
			return fmt.Errorf("the URL %q is not valid: %v", manifest, err)
		}
		builder.URL(3, url)
	default:
		manifest := bytes.NewBufferString(manifest)
		builder.Stream(manifest, "manifest")
	}

	infos, err := builder.Do().Infos()

	if err != nil {
		return err
	}

	for _, info := range infos {
		fmt.Printf("Creating %s: %s", info.Object.GetObjectKind().GroupVersionKind().Kind, info.Name)
		if info.Namespace != "" {
			fmt.Printf(" in namespace: %s", info.Namespace)
		}
		fmt.Println()

		obj, err := resource.NewHelper(info.Client, info.Mapping).Create(info.Namespace, true, info.Object)
		if err != nil {
			fmt.Printf("Warning: failed to create resource: %s\n", err)
		}
		_ = obj
	}

	return nil
}

func DeleteResources(kubeContext, manifestYaml string) error {
	getter := genericclioptions.NewConfigFlags(true)
	getter.Context = &kubeContext
	factory := cmdutil.NewFactory(getter)

	manifest := bytes.NewBufferString(manifestYaml)
	infos, err := factory.NewBuilder().Unstructured().Stream(manifest, "test").Do().Infos()

	if err != nil {
		return err
	}

	for _, info := range infos {
		fmt.Printf("Deleting %s: %s", info.Object.GetObjectKind().GroupVersionKind().Kind, info.Name)
		if info.Namespace != "" {
			fmt.Printf(" in namespace: %s", info.Namespace)
		}
		fmt.Println()

		obj, err := resource.NewHelper(info.Client, info.Mapping).Delete(info.Namespace, info.Name)
		if err != nil {
			fmt.Printf("Warning: failed to delete resource: %s\n", err)
		}
		_ = obj
	}

	return nil
}
