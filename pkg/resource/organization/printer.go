package organization

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"

	"github.com/aws/aws-sdk-go/service/organizations"
)

type OrganizationPrinter struct {
	Organization *organizations.Organization
}

func NewPrinter(Organization *organizations.Organization) *OrganizationPrinter {
	return &OrganizationPrinter{Organization}
}

func (p *OrganizationPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Id", "Feature Set", "Master Account"})

	o := p.Organization

	table.AppendRow([]string{
		aws.StringValue(o.Id),
		aws.StringValue(o.FeatureSet),
		aws.StringValue(o.MasterAccountId),
	})

	table.Print(writer)

	return nil
}

func (p *OrganizationPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.Organization)
}

func (p *OrganizationPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.Organization)
}
