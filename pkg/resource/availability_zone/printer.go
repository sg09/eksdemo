package availability_zone

import (
	"eksdemo/pkg/printer"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type ZonePrinter struct {
	zones []types.AvailabilityZone
}

func NewPrinter(zones []types.AvailabilityZone) *ZonePrinter {
	return &ZonePrinter{zones}
}

func (p *ZonePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Name", "Type", "Opt In Status", "Id"})

	for _, z := range p.zones {
		table.AppendRow([]string{
			aws.ToString(z.ZoneName),
			aws.ToString(z.ZoneType),
			string(z.OptInStatus),
			aws.ToString(z.ZoneId),
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
