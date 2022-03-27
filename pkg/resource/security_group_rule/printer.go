package security_group_rule

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"fmt"
	"io"
	"strconv"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type SecurityGroupRulePrinter struct {
	securityGroupRules []*ec2.SecurityGroupRule
	clusterName        string
}

func NewPrinter(securityGroupRules []*ec2.SecurityGroupRule, clusterName string) *SecurityGroupRulePrinter {
	return &SecurityGroupRulePrinter{securityGroupRules, clusterName}
}

func (p *SecurityGroupRulePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Rule Id", "Proto", "Ports", "Source", "Description"})

	for _, sgr := range p.securityGroupRules {
		fromPort := aws.Int64Value(sgr.FromPort)
		toPort := aws.Int64Value(sgr.ToPort)

		id := aws.StringValue(sgr.SecurityGroupRuleId)
		if aws.BoolValue(sgr.IsEgress) {
			id = "*" + id
		}

		ports := "All"
		if fromPort != -1 {
			if fromPort == toPort {
				ports = strconv.FormatInt(fromPort, 10)
			} else {
				ports = strconv.FormatInt(fromPort, 10) + "-" + strconv.FormatInt(toPort, 10)
			}
		}

		protocol := "All"
		if aws.StringValue(sgr.IpProtocol) != "-1" {
			protocol = aws.StringValue(sgr.IpProtocol)
		}

		source := "-"
		if sgr.ReferencedGroupInfo != nil {
			source = aws.StringValue(sgr.ReferencedGroupInfo.GroupId)
		}
		if sgr.CidrIpv4 != nil {
			source = aws.StringValue(sgr.CidrIpv4)
		}

		table.AppendRow([]string{
			id,
			protocol,
			ports,
			source,
			aws.StringValue(sgr.Description),
		})
	}

	table.SeparateRows()
	table.Print(writer)
	if len(p.securityGroupRules) > 0 {
		fmt.Println("* Indicates egress rule")
	}

	return nil
}

func (p *SecurityGroupRulePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.securityGroupRules)
}

func (p *SecurityGroupRulePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.securityGroupRules)
}
