package target_group

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

type Getter struct{}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var vpcId string

	cluster := options.Common().Cluster
	if cluster != nil {
		vpcId = aws.StringValue(cluster.ResourcesVpcConfig.VpcId)
	}

	targetGroups, err := aws.ELBDescribeTargetGroups(name)
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

func (g *Getter) GetTargetGroupByName(name string) (*elbv2.TargetGroup, error) {
	tg, err := aws.ELBDescribeTargetGroups(name)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case elbv2.ErrCodeTargetGroupNotFoundException:
				return nil, fmt.Errorf("target-group %q not found", name)
			}
		}
		return nil, err
	}

	return tg[0], nil
}
