package amp

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/amp/types"
	"github.com/hako/durafmt"
)

type AmpPrinter struct {
	Workspaces []*types.WorkspaceDescription
}

func NewPrinter(Workspaces []*types.WorkspaceDescription) *AmpPrinter {
	return &AmpPrinter{Workspaces}
}

func (p *AmpPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Alias", "Workspace Id"})

	for _, w := range p.Workspaces {

		age := durafmt.ParseShort(time.Since(aws.TimeValue(w.CreatedAt)))

		table.AppendRow([]string{
			age.String(),
			string(w.Status.StatusCode),
			awssdk.ToString(w.Alias),
			awssdk.ToString(w.WorkspaceId),
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
