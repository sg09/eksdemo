package volume

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
	volumes, err := aws.EC2DescribeVolumes(id)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(volumes))
}
