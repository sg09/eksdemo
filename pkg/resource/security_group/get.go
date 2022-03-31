package security_group

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/load_balancer"
	"eksdemo/pkg/resource/network_interface"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type Getter struct {
	eniGetter network_interface.Getter
	elbGetter load_balancer.Getter
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	sgOptions, ok := options.(*SecurityGroupOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to SecurityGroupOptions")
	}

	var err error
	var securityGroups []*ec2.SecurityGroup

	if sgOptions.NetworkInterfaceId != "" {
		securityGroups, err = g.GetSecurityGroupsByNetworkInterface(sgOptions.NetworkInterfaceId)
	} else if sgOptions.LoadBalancerName != "" {
		securityGroups, err = g.GetSecurityGroupsByLoadBalancerName(sgOptions.LoadBalancerName)
	} else {
		vpcId := ""

		cluster := options.Common().Cluster
		if cluster != nil {
			vpcId = aws.StringValue(cluster.ResourcesVpcConfig.VpcId)
		}

		securityGroups, err = g.GetSecurityGroupsByIdAndVpcFilter(id, vpcId)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(securityGroups, options.Common().ClusterName))
}

func (g *Getter) GetSecurityGroupsByIdAndVpcFilter(id, vpcId string) ([]*ec2.SecurityGroup, error) {
	return aws.EC2DescribeSecurityGroups(id, vpcId, []string{})
}

func (g *Getter) GetSecurityGroupsByLoadBalancerName(name string) ([]*ec2.SecurityGroup, error) {
	sgIds, err := g.elbGetter.GetSecurityGroupIdsForLoadBalancer(name)
	if err != nil {
		return nil, err
	}

	return aws.EC2DescribeSecurityGroups("", "", sgIds)
}

func (g *Getter) GetSecurityGroupsByNetworkInterface(networkInterfaceId string) ([]*ec2.SecurityGroup, error) {
	networkInterface, err := g.eniGetter.GetNetworkInterfaceById(networkInterfaceId)
	if err != nil {
		return nil, err
	}

	if networkInterface == nil {
		return nil, nil
	}

	securityGroupIds := []string{}
	for _, groupIdentifier := range networkInterface.Groups {
		securityGroupIds = append(securityGroupIds, aws.StringValue(groupIdentifier.GroupId))
	}

	if len(securityGroupIds) == 0 {
		return nil, nil
	}

	return aws.EC2DescribeSecurityGroups("", "", securityGroupIds)
}
