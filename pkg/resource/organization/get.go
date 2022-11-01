package organization

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"os"
)

type Getter struct {
	organizationsClient *aws.OrganizationsClient
}

func NewGetter(organizationsClient *aws.OrganizationsClient) *Getter {
	return &Getter{organizationsClient}
}

func (g *Getter) Init() {
	if g.organizationsClient == nil {
		g.organizationsClient = aws.NewOrganizationsClient()
	}
}

func (g *Getter) Get(alias string, output printer.Output, options resource.Options) error {
	org, err := g.organizationsClient.DescribeOrganization()
	if err != nil {
		return aws.FormatErrorAsMessageOnly(err)
	}

	return output.Print(os.Stdout, NewPrinter(org))
}
