package ssm_node

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"os"
)

type Getter struct {
	resource.EmptyInit
}

func (g *Getter) Get(instanceId string, output printer.Output, options resource.Options) error {
	nodes, err := aws.NewSSMClient().DescribeInstanceInformation(instanceId)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(nodes))
}
