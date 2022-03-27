package certificate

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/hako/durafmt"
)

type CertificatePrinter struct {
	certs []*acm.CertificateDetail
}

func NewPrinter(certs []*acm.CertificateDetail) *CertificatePrinter {
	return &CertificatePrinter{certs}
}

func (p *CertificatePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Name", "Status", "In Use"})

	for _, c := range p.certs {
		age := durafmt.ParseShort(time.Since(aws.TimeValue(c.CreatedAt)))

		var inUse string
		if len(c.InUseBy) > 0 {
			inUse = "Yes"
		} else {
			inUse = "No"
		}

		table.AppendRow([]string{
			age.String(),
			aws.StringValue(c.DomainName),
			aws.StringValue(c.Status),
			inUse,
		})
	}

	table.Print(writer)

	return nil
}

func (p *CertificatePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.certs)
}

func (p *CertificatePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.certs)
}
