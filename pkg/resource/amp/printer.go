package amp

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/service/prometheusservice"
	"github.com/hako/durafmt"
)

type AmpPrinter struct {
	Workspaces []*prometheusservice.WorkspaceDescription
}

func NewPrinter(Workspaces []*prometheusservice.WorkspaceDescription) *AmpPrinter {
	return &AmpPrinter{Workspaces}
}

func (p *AmpPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Alias", "Workspace Id"})

	for _, w := range p.Workspaces {

		age := durafmt.ParseShort(time.Since(aws.TimeValue(w.CreatedAt)))

		table.AppendRow([]string{
			age.String(),
			aws.StringValue(w.Status.StatusCode),
			aws.StringValue(w.Alias),
			aws.StringValue(w.WorkspaceId),
		})
	}

	table.Print(writer)

	return nil
}

func (p *AmpPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.Workspaces)
}

func (p *AmpPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.Workspaces)
}
