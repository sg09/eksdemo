package security_group

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/load_balancer"
	"eksdemo/pkg/resource/network_interface"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type Getter struct {
	ec2Client *aws.EC2Client
	eniGetter *network_interface.Getter
	elbGetter *load_balancer.Getter
}

func NewGetter(ec2Client *aws.EC2Client) *Getter {
	return &Getter{ec2Client, network_interface.NewGetter(ec2Client), load_balancer.NewGetter()}
}

func (g *Getter) Init() {
	if g.ec2Client == nil {
		g.ec2Client = aws.NewEC2Client()
	}
	g.eniGetter = network_interface.NewGetter(g.ec2Client)
	g.elbGetter = load_balancer.NewGetter()
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	sgOptions, ok := options.(*SecurityGroupOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to SecurityGroupOptions")
	}

	var err error
	var securityGroups []types.SecurityGroup

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

	return output.Print(os.Stdout, NewPrinter(securityGroups))
}

func (g *Getter) GetSecurityGroupsByIdAndVpcFilter(id, vpcId string) ([]types.SecurityGroup, error) {
	return g.ec2Client.DescribeSecurityGroups(id, vpcId, []string{})
}

func (g *Getter) GetSecurityGroupsByLoadBalancerName(name string) ([]types.SecurityGroup, error) {
	sgIds, err := g.elbGetter.GetSecurityGroupIdsForLoadBalancer(name)
	if err != nil {
		return nil, err
	}

	return g.ec2Client.DescribeSecurityGroups("", "", sgIds)
}

func (g *Getter) GetSecurityGroupsByNetworkInterface(networkInterfaceId string) ([]types.SecurityGroup, error) {
	networkInterface, err := g.eniGetter.GetNetworkInterfaceById(networkInterfaceId)
	if err != nil {
		return nil, err
	}

	securityGroupIds := []string{}
	for _, groupIdentifier := range networkInterface.Groups {
		securityGroupIds = append(securityGroupIds, aws.StringValue(groupIdentifier.GroupId))
	}

	if len(securityGroupIds) == 0 {
		return nil, nil
	}

	return g.ec2Client.DescribeSecurityGroups("", "", securityGroupIds)
}
