package dns_record

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"
)

const MAX_NAME_WITH_UNDERSCORE_LENGTH int = 17
const MAX_RECORD_LENGTH int = 70

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
		name := aws.StringValue(rs.Name)

		if strings.HasPrefix(name, "_") {
			if len(name) > MAX_NAME_WITH_UNDERSCORE_LENGTH {
				name = name[:MAX_NAME_WITH_UNDERSCORE_LENGTH] + "..."
			}
		} else {
			name = strings.TrimSuffix(name, ".")
		}

		records := ""
		if rs.AliasTarget != nil {
			records = aws.StringValue(rs.AliasTarget.DNSName)
		} else {
			for i, rec := range rs.ResourceRecords {
				if i == 0 {
					records = aws.StringValue(rec.Value)
					if len(records) > MAX_RECORD_LENGTH {
						records = records[:MAX_RECORD_LENGTH] + "..."
					}
				} else {
					records += "\n" + aws.StringValue(rec.Value)
				}
			}
		}

		table.AppendRow([]string{
			name,
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
