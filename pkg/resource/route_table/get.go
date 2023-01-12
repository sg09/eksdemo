package route_table

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
	rtOptions, ok := options.(*RouteTableOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to RouteTableOptions")
	}

	cluster := options.Common().Cluster
	filters := []types.Filter{}

	if id != "" {
		filters = append(filters, aws.NewEC2RouteTableFilter(id))
	}

	if cluster != nil {
		filters = append(filters, aws.NewEC2VpcFilter(awssdk.ToString(cluster.ResourcesVpcConfig.VpcId)))
	}

	if rtOptions.VpcId != "" {
		filters = append(filters, aws.NewEC2VpcFilter(rtOptions.VpcId))
	}

	routeTables, err := g.ec2Client.DescribeRouteTables(filters)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(routeTables))
}
