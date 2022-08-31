package availability_zone

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type ZonePrinter struct {
	zones []*ec2.AvailabilityZone
}

func NewPrinter(zones []*ec2.AvailabilityZone) *ZonePrinter {
	return &ZonePrinter{zones}
}

func (p *ZonePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Name", "Type", "Opt In Status", "Id"})

	for _, z := range p.zones {
		table.AppendRow([]string{
			aws.StringValue(z.ZoneName),
			aws.StringValue(z.ZoneType),
			aws.StringValue(z.OptInStatus),
			aws.StringValue(z.ZoneId),
		})
	}
	table.Print(writer)

	return nil
}

func (p *ZonePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.zones)
}

func (p *ZonePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.zones)
}
