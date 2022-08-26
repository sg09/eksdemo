package acm_certificate

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"regexp"
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
	table.SetHeader([]string{"Age", "Id", "Name", "Status", "In Use"})

	resourceId := regexp.MustCompile(`[^:/]*$`)

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
			resourceId.FindString(aws.StringValue(c.CertificateArn)),
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
