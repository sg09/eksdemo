package network_interface

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/load_balancer"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type Getter struct {
	resource.EmptyInit
	elbGetter load_balancer.Getter
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
		vpcId = aws.StringValue(cluster.ResourcesVpcConfig.VpcId)
	}

	if eniOptions.LoadBalancerName != "" {
		elbs, err := g.elbGetter.GetLoadBalancers(eniOptions.LoadBalancerName)
		if err != nil {
			return err
		}

		// Identify the ENIs for a LoadBalancer using Description as described below
		// https://aws.amazon.com/premiumsupport/knowledge-center/elb-find-load-balancer-IP/
		if len(elbs.V1) > 0 {
			description = "ELB " + *elbs.V1[0].LoadBalancerName
		} else if len(elbs.V2) > 0 {
			elb := elbs.V2[0]
			lbName := aws.StringValue(elb.LoadBalancerName)
			lbArn := aws.StringValue(elb.LoadBalancerArn)
			lbId := lbArn[strings.LastIndex(lbArn, "/")+1:]

			switch aws.StringValue(elb.Type) {
			case "application":
				description = fmt.Sprintf("ELB app/%s/%s", lbName, lbId)
			case "network":
				description = fmt.Sprintf("ELB net/%s/%s", lbName, lbId)
			default:
				return fmt.Errorf("load balancer type %q not supported", aws.StringValue(elb.Type))
			}
		}
	}

	networkInterfaces, err := aws.EC2DescribeNetworkInterfaces(
		id, vpcId, description, eniOptions.InstanceId, eniOptions.IpAddress, eniOptions.SecurityGroupId,
	)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(networkInterfaces))
}

func (g *Getter) GetNetworkInterfaceById(id string) (*ec2.NetworkInterface, error) {
	networkInterfaces, err := aws.EC2DescribeNetworkInterfaces(id, "", "", "", "", "")
	if err != nil {
		return nil, err
	}

	if len(networkInterfaces) == 0 {
		return nil, resource.NotFoundError(fmt.Sprintf("Elastic Network Interface %q not found", id))
	}

	return networkInterfaces[0], nil
}
