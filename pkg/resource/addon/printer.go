package addon

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"time"

	"github.com/hako/durafmt"

	"github.com/aws/aws-sdk-go/service/eks"
)

type AddonPrinter struct {
	addons []*eks.Addon
}

func NewPrinter(clusters []*eks.Addon) *AddonPrinter {
	return &AddonPrinter{clusters}
}

func (p *AddonPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Name", "Version"})

	for _, addon := range p.addons {
		age := durafmt.ParseShort(time.Since(*addon.CreatedAt))
		name := aws.StringValue(addon.AddonName)

		table.AppendRow([]string{
			age.String(),
			aws.StringValue(addon.Status),
			name,
			aws.StringValue(addon.AddonVersion),
		})
	}

	table.Print(writer)

	return nil
}

func (p *AddonPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.addons)
}

func (p *AddonPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.addons)
}
