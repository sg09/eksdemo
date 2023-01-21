package amp

import (
	"eksdemo/pkg/printer"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hako/durafmt"
)

type AmpPrinter struct {
	Workspaces []AmpWorkspace
}

func NewPrinter(workspaces []AmpWorkspace) *AmpPrinter {
	return &AmpPrinter{workspaces}
}

func (p *AmpPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Alias", "Workspace Id"})

	for _, w := range p.Workspaces {

		age := durafmt.ParseShort(time.Since(aws.ToTime(w.Workspace.CreatedAt)))

		table.AppendRow([]string{
			age.String(),
			string(w.Workspace.Status.StatusCode),
			aws.ToString(w.Workspace.Alias),
			aws.ToString(w.Workspace.WorkspaceId),
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
