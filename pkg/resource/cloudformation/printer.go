package cloudformation

import (
	"eksdemo/pkg/printer"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/hako/durafmt"
)

type CloudFormationPrinter struct {
	Stacks []types.Stack
}

func NewPrinter(stacks []types.Stack) *CloudFormationPrinter {
	return &CloudFormationPrinter{stacks}
}

func (p *CloudFormationPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Name"})

	for _, s := range p.Stacks {

		age := durafmt.ParseShort(time.Since(aws.ToTime(s.CreationTime)))

		table.AppendRow([]string{
			age.String(),
			string(s.StackStatus),
			aws.ToString(s.StackName),
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
