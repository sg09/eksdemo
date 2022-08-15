package ec2_instance

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
	var vpcId string

	cluster := options.Common().Cluster
	if cluster != nil {
		vpcId = aws.StringValue(cluster.ResourcesVpcConfig.VpcId)
	}

	reservations, err := aws.EC2DescribeInstances(id, vpcId)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(reservations))
}

func (g *Getter) GetInstanceById(id string) (*ec2.Instance, error) {
	reservations, err := aws.EC2DescribeInstances(id, "")
	if err != nil {
		return nil, err
	}

	if len(reservations) == 0 {
		return nil, fmt.Errorf("ec2-instance %q not found", id)
	}

	return reservations[0].Instances[0], nil
}
