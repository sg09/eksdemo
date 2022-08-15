package nat_gateway

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hako/durafmt"
)

const maxNameLength = 35

type NatGatewayPrinter struct {
	nats []*ec2.NatGateway
}

func NewPrinter(nats []*ec2.NatGateway) *NatGatewayPrinter {
	return &NatGatewayPrinter{nats}
}

func (p *NatGatewayPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "State", "Id", "Name"})

	for _, n := range p.nats {
		age := durafmt.ParseShort(time.Since(aws.TimeValue(n.CreateTime)))

		table.AppendRow([]string{
			age.String(),
			aws.StringValue(n.State),
			aws.StringValue(n.NatGatewayId),
			p.getNatGatewayName(n),
		})
	}
	table.Print(writer)

	return nil
}

func (p *NatGatewayPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.nats)
}

func (p *NatGatewayPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.nats)
}

func (p *NatGatewayPrinter) getNatGatewayName(nat *ec2.NatGateway) string {
	name := ""
	for _, tag := range nat.Tags {
		if aws.StringValue(tag.Key) == "Name" {
			name = aws.StringValue(tag.Value)

			if len(name) > maxNameLength {
				name = name[:maxNameLength-3] + "..."
			}
			continue
		}
	}
	return name
}
