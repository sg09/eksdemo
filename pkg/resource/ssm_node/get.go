package ssm_node

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"os"
)

type Getter struct{}

func (g *Getter) Get(instanceId string, output printer.Output, options resource.Options) error {
	vpcs, err := aws.SSMDescribeInstanceInformation(instanceId)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(vpcs))
}
