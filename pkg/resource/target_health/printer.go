package target_health

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"strconv"

	"github.com/aws/aws-sdk-go/service/elbv2"
)

type TargetHealthPrinter struct {
	targets []*elbv2.TargetHealthDescription
}

func NewPrinter(targets []*elbv2.TargetHealthDescription) *TargetHealthPrinter {
	return &TargetHealthPrinter{targets}
}

func (p *TargetHealthPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"State", "Id", "Port", "Details"})

	for _, t := range p.targets {
		table.AppendRow([]string{
			aws.StringValue(t.TargetHealth.State),
			aws.StringValue(t.Target.Id),
			strconv.FormatInt(aws.Int64Value(t.Target.Port), 10),
			aws.StringValue(t.TargetHealth.Description),
		})
	}
	table.Print(writer)

	return nil
}

func (p *TargetHealthPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.targets)
}

func (p *TargetHealthPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.targets)
}
