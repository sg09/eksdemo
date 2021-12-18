package dns_record

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"
)

type RecordSetPrinter struct {
	recordSets []*route53.ResourceRecordSet
}

func NewPrinter(recordSets []*route53.ResourceRecordSet) *RecordSetPrinter {
	return &RecordSetPrinter{recordSets}
}

func (p *RecordSetPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Name", "Type", "Value"})

	for _, rs := range p.recordSets {
		recordType := aws.StringValue(rs.Type)
		if recordType == "SOA" || recordType == "TXT" {
			continue
		}
		name := aws.StringValue(rs.Name)
		if strings.HasPrefix(name, "_") {
			continue
		}

		records := ""
		if rs.AliasTarget != nil {
			records = aws.StringValue(rs.AliasTarget.DNSName)
		} else {
			for i, rec := range rs.ResourceRecords {
				if i == 0 {
					records = aws.StringValue(rec.Value)
				} else {
					records += "\n" + aws.StringValue(rec.Value)
				}
			}
		}

		table.AppendRow([]string{
			strings.TrimSuffix(name, "."),
			recordType,
			records,
		})
	}

	table.Print(writer)

	return nil
}

func (p *RecordSetPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.recordSets)
}

func (p *RecordSetPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.recordSets)
}
