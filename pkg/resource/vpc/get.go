package vpc

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"os"
)

type Getter struct {
	resource.EmptyInit
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	var vpcId string

	cluster := options.Common().Cluster
	if cluster != nil {
		vpcId = aws.StringValue(cluster.ResourcesVpcConfig.VpcId)
	}

	vpcs, err := aws.EC2DescribeVpcs(id, vpcId)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(vpcs))
}
