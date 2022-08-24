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
	vpcs          []*ec2.Vpc
	multipleCidrs bool
}

func NewPrinter(vpcs []*ec2.Vpc) *VpcPrinter {
	return &VpcPrinter{vpcs: vpcs}
}

func (p *VpcPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Id", "Name", "IPv4 CIDR(s)", "IPv6 CIDR(s)"})

	for _, vpc := range p.vpcs {
		name := p.getVpcName(vpc)
		if aws.BoolValue(vpc.IsDefault) {
			name += "*"
		}

		vpcCidr := aws.StringValue(vpc.CidrBlock)
		v4Cidrs := []string{vpcCidr}

		for _, cbas := range vpc.CidrBlockAssociationSet {
			cbasCidr := aws.StringValue(cbas.CidrBlock)
			if cbasCidr != vpcCidr && aws.StringValue(cbas.CidrBlockState.State) == "associated" {
				v4Cidrs = append(v4Cidrs, cbasCidr)
			}
		}

		v6Cidrs := make([]string, 0, len(vpc.Ipv6CidrBlockAssociationSet))
		for _, cba := range vpc.Ipv6CidrBlockAssociationSet {
			v6Cidrs = append(v6Cidrs, aws.StringValue(cba.Ipv6CidrBlock))
		}

		if len(v4Cidrs) > 1 || len(v6Cidrs) > 1 {
			p.multipleCidrs = true
		}

		table.AppendRow([]string{
			aws.StringValue(vpc.VpcId),
			name,
			strings.Join(v4Cidrs, "\n"),
			strings.Join(v6Cidrs, "\n"),
		})
	}

	if p.multipleCidrs {
		table.SeparateRows()
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
	for _, tag := range vpc.Tags {
		if aws.StringValue(tag.Key) == "Name" {
			return aws.StringValue(tag.Value)
		}
	}
	return ""
}
