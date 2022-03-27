package network_interface

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type Getter struct{}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	eniOptions, ok := options.(*NetworkInterfaceOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to NetworkInterfaceOptions")
	}

	vpcId := ""

	cluster := options.Common().Cluster
	if cluster != nil {
		vpcId = aws.StringValue(cluster.ResourcesVpcConfig.VpcId)
	}

	networkInterfaces, err := aws.EC2DescribeNetworkInterfaces(
		id, vpcId, eniOptions.InstanceId, eniOptions.IpAddress, eniOptions.SecurityGroupId,
	)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(networkInterfaces, options.Common().ClusterName))
}

func (g *Getter) GetNetworkInterfaceById(id string) (*ec2.NetworkInterface, error) {
	networkInterfaces, err := aws.EC2DescribeNetworkInterfaces(id, "", "", "", "")
	if err != nil {
		return nil, err
	}

	if len(networkInterfaces) == 0 {
		return nil, resource.NotFoundError(fmt.Sprintf("Elastic Network Interface %q not found", id))
	}

	return networkInterfaces[0], nil
}
