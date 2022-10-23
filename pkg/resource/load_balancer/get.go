package load_balancer

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

type LoadBalancers struct {
	V1 []*elb.LoadBalancerDescription
	V2 []*elbv2.LoadBalancer
}

type Getter struct {
	resource.EmptyInit
	elbs *LoadBalancers
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) (err error) {
	g.elbs, err = g.GetLoadBalancers(name)
	if err != nil {
		return err
	}

	cluster := options.Common().Cluster
	if cluster != nil {
		g.filterByVpc(aws.StringValue(cluster.ResourcesVpcConfig.VpcId))
	}

	return output.Print(os.Stdout, NewPrinter(g.elbs))
}

func (g *Getter) GetLoadBalancers(name string) (elbs *LoadBalancers, err error) {
	elbs = &LoadBalancers{}

	elbs.V1, err = aws.ELBDescribeLoadBalancersv1(name)
	if err != nil {
		// Return all errors except NotFound
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case elb.ErrCodeAccessPointNotFoundException:
				break
			default:
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	elbs.V2, err = aws.ELBDescribeLoadBalancersv2(name)
	if err != nil {
		// Return all errors except NotFound
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case elbv2.ErrCodeLoadBalancerNotFoundException:
				break
			default:
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if name != "" && len(elbs.V1) == 0 && len(elbs.V2) == 0 {
		return nil, fmt.Errorf("load balancer %q not found", name)
	}

	return elbs, nil
}

func (g *Getter) GetSecurityGroupIdsForLoadBalancer(name string) (ids []string, err error) {
	g.elbs, err = g.GetLoadBalancers(name)
	if err != nil {
		return nil, err
	}

	// Check for the unlikely but possible scenario with elbv1 and elbv2 with same name
	if len(g.elbs.V1) > 0 && len(g.elbs.V2) > 0 {
		return nil, fmt.Errorf("multiple load balancers with name %q", name)
	}

	if len(g.elbs.V2) > 0 {
		return aws.StringValueSlice(g.elbs.V2[0].SecurityGroups), nil
	}

	if len(g.elbs.V1) > 0 {
		return aws.StringValueSlice(g.elbs.V1[0].SecurityGroups), nil
	}

	return nil, nil
}

func (g *Getter) filterByVpc(id string) {
	filteredV1 := make([]*elb.LoadBalancerDescription, 0, len(g.elbs.V1))
	filteredV2 := make([]*elbv2.LoadBalancer, 0, len(g.elbs.V2))

	for _, v1 := range g.elbs.V1 {
		if aws.StringValue(v1.VPCId) == id {
			filteredV1 = append(filteredV1, v1)
		}
	}

	for _, v2 := range g.elbs.V2 {
		if aws.StringValue(v2.VpcId) == id {
			filteredV2 = append(filteredV2, v2)
		}
	}

	g.elbs.V1 = filteredV1
	g.elbs.V2 = filteredV2
}
