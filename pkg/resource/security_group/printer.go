package security_group

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/eksctl"
	"eksdemo/pkg/printer"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type SecurityGroupPrinter struct {
	securityGroups []*ec2.SecurityGroup
	clusterName    string
}

func NewPrinter(securityGroups []*ec2.SecurityGroup, clusterName string) *SecurityGroupPrinter {
	return &SecurityGroupPrinter{securityGroups, clusterName}
}

func (p *SecurityGroupPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Group Id", "Description", "Name"})

	namePrefix := eksctl.TagNamePrefix(p.clusterName)
	prefixFound := false

	for _, sg := range p.securityGroups {
		name := p.getName(sg)

		if strings.HasPrefix(name, namePrefix) {
			prefixFound = true
			name = "*" + name[len(namePrefix):]
		}

		table.AppendRow([]string{
			aws.StringValue(sg.GroupId),
			aws.StringValue(sg.Description),
			name,
		})
	}

	table.SeparateRows()
	table.Print(writer)
	if prefixFound {
		fmt.Printf("* Name begins with %q\n", namePrefix)
	}

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
