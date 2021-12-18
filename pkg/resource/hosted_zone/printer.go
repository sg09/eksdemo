package hosted_zone

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"
)

type HostedZonePrinter struct {
	zones []*route53.HostedZone
}

func NewPrinter(zones []*route53.HostedZone) *HostedZonePrinter {
	return &HostedZonePrinter{zones}
}

func (p *HostedZonePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Name", "Type", "Records", "Zone Id"})

	for _, z := range p.zones {
		var zoneType string
		if aws.BoolValue(z.Config.PrivateZone) {
			zoneType = "Private"
		} else {
			zoneType = "Public"
		}

		table.AppendRow([]string{
			strings.TrimSuffix(aws.StringValue(z.Name), "."),
			zoneType,
			strconv.Itoa(int(aws.Int64Value(z.ResourceRecordSetCount))),
			strings.TrimPrefix(aws.StringValue(z.Id), "/hostedzone/"),
		})
	}

	table.Print(writer)

	return nil
}

func (p *HostedZonePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.zones)
}

func (p *HostedZonePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.zones)
}
