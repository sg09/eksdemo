package amg

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"strings"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/grafana/types"
	"github.com/hako/durafmt"
)

type AmgPrinter struct {
	Workspaces []*types.WorkspaceDescription
}

func NewPrinter(Workspaces []*types.WorkspaceDescription) *AmgPrinter {
	return &AmgPrinter{Workspaces}
}

func (p *AmgPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Name", "Workspace Id", "Auth"})

	for _, w := range p.Workspaces {

		age := durafmt.ParseShort(time.Since(aws.TimeValue(w.Created)))

		table.AppendRow([]string{
			age.String(),
			string(w.Status),
			awssdk.ToString(w.Name),
			awssdk.ToString(w.Id),
			strings.Join(toStringSlice(w.Authentication.Providers), ","),
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

func toStringSlice(apt []types.AuthenticationProviderTypes) []string {
	ss := make([]string, len(apt))
	for i, v := range apt {
		ss[i] = string(v)
	}
	return ss
}
