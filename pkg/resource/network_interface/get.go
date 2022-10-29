package network_interface

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/load_balancer"
	"fmt"
	"os"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type Getter struct {
	ec2Client          *aws.EC2Client
	loadBalancerGetter *load_balancer.Getter
}

func NewGetter(ec2Client *aws.EC2Client) *Getter {
	return &Getter{ec2Client, load_balancer.NewGetter(aws.NewElasticloadbalancingClientv1(), aws.NewElasticloadbalancingClientv2())}
}

func (g *Getter) Init() {
	if g.ec2Client == nil {
		g.ec2Client = aws.NewEC2Client()
	}
	g.loadBalancerGetter = load_balancer.NewGetter(aws.NewElasticloadbalancingClientv1(), aws.NewElasticloadbalancingClientv2())
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	eniOptions, ok := options.(*NetworkInterfaceOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to NetworkInterfaceOptions")
	}

	description := ""
	vpcId := ""

	cluster := options.Common().Cluster
	if cluster != nil {
		vpcId = awssdk.ToString(cluster.ResourcesVpcConfig.VpcId)
	}

	if eniOptions.LoadBalancerName != "" {
		elbs, err := g.loadBalancerGetter.GetLoadBalancers(eniOptions.LoadBalancerName)
		if err != nil {
			return err
		}

		// Identify the ENIs for a LoadBalancer using Description as described below
		// https://aws.amazon.com/premiumsupport/knowledge-center/elb-find-load-balancer-IP/
		if len(elbs.V1) > 0 {
			description = "ELB " + *elbs.V1[0].LoadBalancerName
		} else if len(elbs.V2) > 0 {
			elb := elbs.V2[0]
			lbName := awssdk.ToString(elb.LoadBalancerName)
			lbArn := awssdk.ToString(elb.LoadBalancerArn)
			lbId := lbArn[strings.LastIndex(lbArn, "/")+1:]

			switch string(elb.Type) {
			case "application":
				description = fmt.Sprintf("ELB app/%s/%s", lbName, lbId)
			case "network":
				description = fmt.Sprintf("ELB net/%s/%s", lbName, lbId)
			default:
				return fmt.Errorf("load balancer type %q not supported", string(elb.Type))
			}
		}
	}

	networkInterfaces, err := g.ec2Client.DescribeNetworkInterfaces(
		id, vpcId, description, eniOptions.InstanceId, eniOptions.IpAddress, eniOptions.SecurityGroupId,
	)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(networkInterfaces, aws.AccountId()))
}

func (g *Getter) GetNetworkInterfaceById(id string) (types.NetworkInterface, error) {
	networkInterfaces, err := g.ec2Client.DescribeNetworkInterfaces(id, "", "", "", "", "")
	if err != nil {
		return types.NetworkInterface{}, err
	}

	if len(networkInterfaces) == 0 {
		return types.NetworkInterface{}, resource.NotFoundError(fmt.Sprintf("network-interface %q not found", id))
	}

	return networkInterfaces[0], nil
}
