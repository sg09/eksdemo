package subnet

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type SubnetPrinter struct {
	subnets       []*ec2.Subnet
	multipleCidrs bool
}

func NewPrinter(subnets []*ec2.Subnet) *SubnetPrinter {
	return &SubnetPrinter{subnets: subnets}
}

func (p *SubnetPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Id", "Zone", "IPv4 CIDR", "Free", "IPv6 CIDR"})

	for _, subnet := range p.subnets {
		v6Cidrs := make([]string, 0, len(subnet.Ipv6CidrBlockAssociationSet))
		for _, cbas := range subnet.Ipv6CidrBlockAssociationSet {
			v6Cidrs = append(v6Cidrs, aws.StringValue(cbas.Ipv6CidrBlock))
		}

		if len(v6Cidrs) == 0 {
			v6Cidrs = []string{"-"}
		} else {
			p.multipleCidrs = true
		}

		table.AppendRow([]string{
			aws.StringValue(subnet.SubnetId),
			aws.StringValue(subnet.AvailabilityZone),
			aws.StringValue(subnet.CidrBlock),
			strconv.Itoa(int(aws.Int64Value(subnet.AvailableIpAddressCount))),
			strings.Join(v6Cidrs, "\n"),
		})
	}

	if p.multipleCidrs {
		table.SeparateRows()
	}

	table.Print(writer)

	return nil
}

func (p *SubnetPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.subnets)
}

func (p *SubnetPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.subnets)
}
