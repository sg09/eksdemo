package target_group

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/load_balancer"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

type Getter struct {
	elbClientv2        *aws.Elasticloadbalancingv2Client
	loadBalancerGetter *load_balancer.Getter
}

func NewGetter(elbClientv2 *aws.Elasticloadbalancingv2Client) *Getter {
	return &Getter{elbClientv2, load_balancer.NewGetter(aws.NewElasticloadbalancingClientv1(), aws.NewElasticloadbalancingClientv2())}
}

func (g *Getter) Init() {
	if g.elbClientv2 == nil {
		g.elbClientv2 = aws.NewElasticloadbalancingClientv2()
	}
	g.loadBalancerGetter = load_balancer.NewGetter(aws.NewElasticloadbalancingClientv1(), aws.NewElasticloadbalancingClientv2())
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
		elbs, err := g.loadBalancerGetter.GetLoadBalancers(tgOptions.LoadBalancerName)
		if err != nil {
			return err
		}
		if len(elbs.V1) > 0 {
			return fmt.Errorf("%q is a classic load balancer", tgOptions.LoadBalancerName)
		}

		lbArn = aws.StringValue(elbs.V2[0].LoadBalancerArn)
	}

	targetGroups, err := g.elbClientv2.DescribeTargetGroups(name, lbArn)
	if err != nil {
		return err
	}

	if vpcId != "" {
		filtered := make([]types.TargetGroup, 0, len(targetGroups))

		for _, tg := range targetGroups {
			if aws.StringValue(tg.VpcId) == vpcId {
				filtered = append(filtered, tg)
			}
		}
		targetGroups = filtered
	}

	return output.Print(os.Stdout, NewPrinter(targetGroups))
}

func (g *Getter) GetTargetGroupByName(name string) (types.TargetGroup, error) {
	tg, err := g.elbClientv2.DescribeTargetGroups(name, "")

	if err != nil {
		var tgnfe *types.TargetGroupNotFoundException
		if errors.As(err, &tgnfe) {
			return types.TargetGroup{}, fmt.Errorf("target-group %q not found", name)
		}
		return types.TargetGroup{}, err
	}

	return tg[0], nil
}
