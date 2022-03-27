package subnet

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/eksctl"
	"eksdemo/pkg/printer"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type SubnetPrinter struct {
	subnets     []*ec2.Subnet
	clusterName string
}

func NewPrinter(subnets []*ec2.Subnet, clusterName string) *SubnetPrinter {
	return &SubnetPrinter{subnets, clusterName}
}

func (p *SubnetPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Id", "Name", "CIDR Block", "Avail IP"})

	namePrefix := eksctl.TagNamePrefix(p.clusterName)
	prefixFound := false

	for _, subnet := range p.subnets {
		name := p.getName(subnet)

		if strings.HasPrefix(name, namePrefix) {
			prefixFound = true
			name = "*" + name[len(namePrefix):]
		}

		table.AppendRow([]string{
			aws.StringValue(subnet.SubnetId),
			name,
			aws.StringValue(subnet.CidrBlock),
			strconv.Itoa(int(aws.Int64Value(subnet.AvailableIpAddressCount))),
		})
	}

	table.Print(writer)
	if prefixFound {
		fmt.Printf("* Name begins with %q\n", namePrefix)
	}

	return nil
}

func (p *SubnetPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.subnets)
}

func (p *SubnetPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.subnets)
}

func (p *SubnetPrinter) getName(subnet *ec2.Subnet) string {
	name := ""
	for _, tag := range subnet.Tags {
		if aws.StringValue(tag.Key) == "Name" {
			name = aws.StringValue(tag.Value)
			continue
		}
	}
	return name
}
