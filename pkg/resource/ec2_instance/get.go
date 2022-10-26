package ec2_instance

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"os"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type Getter struct {
	ec2Client *aws.EC2Client
}

func NewGetter(ec2Client *aws.EC2Client) *Getter {
	return &Getter{ec2Client}
}

func (g *Getter) Init() {
	if g.ec2Client == nil {
		g.ec2Client = aws.NewEC2Client()
	}
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	cluster := options.Common().Cluster
	filters := []types.Filter{}

	if id != "" {
		filters = append(filters, aws.NewEC2InstanceFilter(id))
	}

	if cluster != nil {
		filters = append(filters, aws.NewEC2VpcFilter(awssdk.ToString(cluster.ResourcesVpcConfig.VpcId)))
	}

	reservations, err := g.ec2Client.DescribeInstances(filters)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(reservations))
}

func (g *Getter) GetInstanceById(id string) (types.Instance, error) {
	reservations, err := g.ec2Client.DescribeInstances([]types.Filter{aws.NewEC2InstanceFilter(id)})
	if err != nil {
		return types.Instance{}, err
	}

	if len(reservations) == 0 {
		return types.Instance{}, fmt.Errorf("ec2-instance %q not found", id)
	}

	return reservations[0].Instances[0], nil
}
