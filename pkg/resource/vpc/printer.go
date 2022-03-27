package vpc

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type VpcPrinter struct {
	vpcs []*ec2.Vpc
}

func NewPrinter(vpcs []*ec2.Vpc) *VpcPrinter {
	return &VpcPrinter{vpcs}
}

func (p *VpcPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Id", "Name", "IPv4 CIDR", "Secondary CIDRs"})

	for _, vpc := range p.vpcs {
		name := p.getVpcName(vpc)
		if aws.BoolValue(vpc.IsDefault) {
			name += "*"
		}

		secondaryCidrs := make([]string, 0, len(vpc.CidrBlockAssociationSet))
		for _, cba := range vpc.CidrBlockAssociationSet {
			if aws.StringValue(cba.CidrBlock) != aws.StringValue(vpc.CidrBlock) &&
				aws.StringValue(cba.CidrBlockState.State) == "associated" {
				secondaryCidrs = append(secondaryCidrs, aws.StringValue(cba.CidrBlock))
			}
		}

		table.AppendRow([]string{
			aws.StringValue(vpc.VpcId),
			name,
			aws.StringValue(vpc.CidrBlock),
			strings.Join(secondaryCidrs, ", "),
		})
	}

	table.Print(writer)
	if len(p.vpcs) > 0 {
		fmt.Println("* Indicates default VPC")
	}

	return nil
}

func (p *VpcPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.vpcs)
}

func (p *VpcPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.vpcs)
}

func (p *VpcPrinter) getVpcName(vpc *ec2.Vpc) string {
	name := ""
	for _, tag := range vpc.Tags {
		if aws.StringValue(tag.Key) == "Name" {
			name = aws.StringValue(tag.Value)
			continue
		}
	}
	return name
}
