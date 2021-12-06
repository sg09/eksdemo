package fargate_profile

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"os"

	"github.com/aws/aws-sdk-go/service/eks"
)

type Getter struct{}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var profiles []*eks.FargateProfile
	var err error

	clusterName := options.Common().ClusterName

	if name != "" {
		profiles, err = g.GetProfileByName(name, clusterName)
	} else {
		profiles, err = g.GetAllProfiles(clusterName)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(profiles))
}

func (g *Getter) GetAllProfiles(clusterName string) ([]*eks.FargateProfile, error) {
	profileNames, err := aws.EksListFargateProfiles(clusterName)
	profiles := make([]*eks.FargateProfile, 0, len(profileNames))

	if err != nil {
		return nil, err
	}

	for _, name := range profileNames {
		result, err := aws.EksDescribeFargateProfile(clusterName, aws.StringValue(name))
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, result)
	}

	return profiles, nil
}

func (g *Getter) GetProfileByName(name, clusterName string) ([]*eks.FargateProfile, error) {
	profile, err := aws.EksDescribeFargateProfile(clusterName, name)
	if err != nil {
		return nil, err
	}

	return []*eks.FargateProfile{profile}, nil
}
