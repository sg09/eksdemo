package amg

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/managedgrafana"
	"github.com/hako/durafmt"
)

type AmgPrinter struct {
	Workspaces []*managedgrafana.WorkspaceDescription
}

func NewPrinter(Workspaces []*managedgrafana.WorkspaceDescription) *AmgPrinter {
	return &AmgPrinter{Workspaces}
}

func (p *AmgPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Name", "Workspace Id", "Auth"})

	for _, w := range p.Workspaces {

		age := durafmt.ParseShort(time.Since(aws.TimeValue(w.Created)))

		table.AppendRow([]string{
			age.String(),
			aws.StringValue(w.Status),
			aws.StringValue(w.Name),
			aws.StringValue(w.Id),
			strings.Join(aws.StringValueSlice(w.Authentication.Providers), ","),
		})
	}

	table.Print(writer)

	return nil
}

func (p *AmgPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.Workspaces)
}

func (p *AmgPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.Workspaces)
}
