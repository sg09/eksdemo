package cloudformation

import (
	"eksdemo/pkg/printer"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/hako/durafmt"
)

type CloudFormationPrinter struct {
	Stacks []*cloudformation.Stack
}

func NewPrinter(Workspaces []*cloudformation.Stack) *CloudFormationPrinter {
	return &CloudFormationPrinter{Workspaces}
}

func (p *CloudFormationPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Name"})

	for _, s := range p.Stacks {

		age := durafmt.ParseShort(time.Since(*s.CreationTime))

		table.AppendRow([]string{
			age.String(),
			*s.StackStatus,
			*s.StackName,
		})
	}

	table.Print(writer)

	return nil
}

func (p *CloudFormationPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.Stacks)
}

func (p *CloudFormationPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.Stacks)
}
