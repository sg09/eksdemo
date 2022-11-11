package network_acl

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
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
	filters := []types.Filter{}

	if cluster := options.Common().Cluster; cluster != nil {
		filters = append(filters, aws.NewEC2VpcFilter(awssdk.ToString(cluster.ResourcesVpcConfig.VpcId)))
	}

	if id != "" {
		filters = append(filters, aws.NewEC2NetworkAclFilter(id))
	}

	nacls, err := g.ec2Client.DescribeNetworkAcls(filters)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(nacls))
}
