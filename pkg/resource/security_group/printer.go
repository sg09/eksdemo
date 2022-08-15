package security_group

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type SecurityGroupPrinter struct {
	securityGroups []*ec2.SecurityGroup
}

func NewPrinter(securityGroups []*ec2.SecurityGroup) *SecurityGroupPrinter {
	return &SecurityGroupPrinter{securityGroups}
}

func (p *SecurityGroupPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Group Id", "Group Name", "Description"})

	for _, sg := range p.securityGroups {
		_ = p.getName(sg)

		table.AppendRow([]string{
			aws.StringValue(sg.GroupId),
			aws.StringValue(sg.GroupName),
			aws.StringValue(sg.Description),
		})
	}

	table.SeparateRows()
	table.Print(writer)

	return nil
}

func (p *SecurityGroupPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.securityGroups)
}

func (p *SecurityGroupPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.securityGroups)
}

func (p *SecurityGroupPrinter) getName(sg *ec2.SecurityGroup) string {
	name := ""
	for _, tag := range sg.Tags {
		if aws.StringValue(tag.Key) == "Name" {
			name = aws.StringValue(tag.Value)
			continue
		}
	}
	return name
}
