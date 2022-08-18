package dns_record

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"
)

const MAX_COMBINED_NAME_AND_RECORD_LENGTH int = 90

type RecordSetPrinter struct {
	recordSets        []*route53.ResourceRecordSet
	longestNameLength int
}

func NewPrinter(recordSets []*route53.ResourceRecordSet) *RecordSetPrinter {
	return &RecordSetPrinter{recordSets: recordSets}
}

func (p *RecordSetPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Name", "Type", "Value"})

	for _, rs := range p.recordSets {
		if l := len(aws.StringValue(rs.Name)); l > p.longestNameLength {
			p.longestNameLength = l
		}
	}

	for _, rs := range p.recordSets {
		records := ""
		if rs.AliasTarget != nil {
			records = p.limitLength(aws.StringValue(rs.AliasTarget.DNSName))
		} else {
			for i, rec := range rs.ResourceRecords {
				if i == 0 {
					records = p.limitLength(aws.StringValue(rec.Value))
				} else {
					records += "\n" + aws.StringValue(rec.Value)
				}
			}
		}

		table.AppendRow([]string{
			strings.TrimSuffix(aws.StringValue(rs.Name), "."),
			aws.StringValue(rs.Type),
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

func (p *RecordSetPrinter) limitLength(record string) string {
	if len(record) > MAX_COMBINED_NAME_AND_RECORD_LENGTH-p.longestNameLength {
		record = record[:MAX_COMBINED_NAME_AND_RECORD_LENGTH-p.longestNameLength-3] + "..."
	}
	return record
}
