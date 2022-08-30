package ssm_node

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hako/durafmt"
)

type NodePrinter struct {
	nodes []*ssm.InstanceInformation
}

func NewPrinter(nodes []*ssm.InstanceInformation) *NodePrinter {
	return &NodePrinter{nodes}
}

func (p *NodePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Ping", "Status", "Instance Id", "IP Address", "Agent Ver", "OS"})

	for _, n := range p.nodes {
		ping := durafmt.ParseShort(time.Since(aws.TimeValue(n.LastPingDateTime)))

		table.AppendRow([]string{
			ping.String(),
			aws.StringValue(n.PingStatus),
			aws.StringValue(n.InstanceId),
			aws.StringValue(n.IPAddress),
			aws.StringValue(n.AgentVersion),
			aws.StringValue(n.PlatformName) + " " + aws.StringValue(n.PlatformVersion),
		})
	}

	table.Print(writer)

	return nil
}

func (p *NodePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.nodes)
}

func (p *NodePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.nodes)
}
