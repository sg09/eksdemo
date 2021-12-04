package subnet

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"os"
)

type Getter struct{}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var vpcId string

	cluster := options.Common().Cluster
	if cluster != nil {
		vpcId = aws.StringValue(cluster.ResourcesVpcConfig.VpcId)
	}

	subnets, err := aws.EC2DescribeSubnets(name, vpcId)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(subnets, options.Common().ClusterName))
}
