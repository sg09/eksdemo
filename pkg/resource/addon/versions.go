package addon

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"os"

	"github.com/aws/aws-sdk-go/service/eks"
)

func NewVersionsResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "addon-versions",
			Description: "EKS Managed Addon Versions",
			Args:        []string{"NAME"},
		},

		Getter: &VersionGetter{},
	}

	res.Options = &resource.CommonOptions{}
	res.CreateFlags = cmd.Flags{}

	return res
}

type VersionGetter struct {
	resource.EmptyInit
}

func (g *VersionGetter) Get(name string, output printer.Output, options resource.Options) error {
	var addonVersions []*eks.AddonInfo
	var err error

	k8sversion := options.Common().KubernetesVersion

	if name != "" {
		addonVersions, err = g.GetAddonVersionsByName(name, k8sversion)
	} else {
		addonVersions, err = g.GetAllAddonVersions(k8sversion)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewVersionPrinter(addonVersions))
}

func (g *VersionGetter) GetAddonVersionsByName(name, k8sversion string) ([]*eks.AddonInfo, error) {
	addonVersions, err := aws.EksDescribeAddonVersions(name, k8sversion)
	if err != nil {
		return nil, err
	}

	return addonVersions, nil
}

func (g *VersionGetter) GetAllAddonVersions(k8sversion string) ([]*eks.AddonInfo, error) {
	addonVersions, err := aws.EksDescribeAddonVersions("", k8sversion)
	if err != nil {
		return nil, err
	}

	return addonVersions, nil
}
