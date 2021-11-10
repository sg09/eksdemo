package organization

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"os"
)

type Getter struct{}

func (g *Getter) Get(alias string, output printer.Output, options resource.Options) error {
	org, err := aws.OrgsDescribeOrganization()
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(org))
}
