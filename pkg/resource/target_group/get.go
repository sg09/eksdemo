package target_group

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/load_balancer"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

type Getter struct {
	elbGetter load_balancer.Getter
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	tgOptions, ok := options.(*TargeGroupOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to TargeGroupOptions")
	}

	var lbArn, vpcId string

	cluster := options.Common().Cluster
	if cluster != nil {
		vpcId = aws.StringValue(cluster.ResourcesVpcConfig.VpcId)
	}

	if tgOptions.LoadBalancerName != "" {
		elbs, err := g.elbGetter.GetLoadBalancers(tgOptions.LoadBalancerName)
		if err != nil {
			return err
		}
		if len(elbs.V1) > 0 {
			return fmt.Errorf("%q is a classic load balancer", tgOptions.LoadBalancerName)
		}

		lbArn = aws.StringValue(elbs.V2[0].LoadBalancerArn)
	}

	targetGroups, err := aws.ELBDescribeTargetGroups(name, lbArn)
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
	tg, err := aws.ELBDescribeTargetGroups(name, "")
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
