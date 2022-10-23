package addon

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"os"

	"github.com/aws/aws-sdk-go/service/eks"
)

type Getter struct {
	resource.EmptyInit
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var addons []*eks.Addon
	var err error

	clusterName := options.Common().ClusterName

	if name != "" {
		addons, err = g.GetAddonsByName(name, clusterName)
	} else {
		addons, err = g.GetAllAddons(clusterName)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(addons))
}

func (g *Getter) GetAddonsByName(name, clusterName string) ([]*eks.Addon, error) {
	addon, err := aws.EksDescribeAddon(clusterName, name)
	if err != nil {
		return nil, err
	}

	return []*eks.Addon{addon}, nil
}

func (g *Getter) GetAllAddons(clusterName string) ([]*eks.Addon, error) {
	addonNames, err := aws.EksListAddons(clusterName)
	addons := make([]*eks.Addon, 0, len(addonNames))

	if err != nil {
		return nil, err
	}

	for _, name := range addonNames {
		result, err := aws.EksDescribeAddon(clusterName, aws.StringValue(name))
		if err != nil {
			return nil, err
		}
		addons = append(addons, result)
	}

	return addons, nil
}
