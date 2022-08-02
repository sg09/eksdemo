package target_group

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"os"

	"github.com/aws/aws-sdk-go/service/elbv2"
)

type Getter struct{}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	var vpcId string

	cluster := options.Common().Cluster
	if cluster != nil {
		vpcId = aws.StringValue(cluster.ResourcesVpcConfig.VpcId)
	}

	targetGroups, err := aws.ELBDescribeTargetGroups(id)
	if err != nil {
		return err
	}

	if vpcId != "" {
		filtered := make([]*elbv2.TargetGroup, 0, len(targetGroups))

		for _, tg := range targetGroups {
			if aws.StringValue(tg.VpcId) == vpcId {
				filtered = append(filtered, tg)
			}
		}
		targetGroups = filtered
	}

	return output.Print(os.Stdout, NewPrinter(targetGroups))
}
