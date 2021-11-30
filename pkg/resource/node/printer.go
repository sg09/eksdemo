package node

import (
	"eksdemo/pkg/printer"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hako/durafmt"
	v1 "k8s.io/api/core/v1"
)

type NodePrinter struct {
	nodes []v1.Node
}

func NewPrinter(nodes []v1.Node) *NodePrinter {
	return &NodePrinter{nodes}
}

func (p *NodePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Name", "Instance Id", "Zone", "Nodegroup", "Type"})

	for _, node := range p.nodes {
		age := durafmt.ParseShort(time.Since(node.CreationTimestamp.Time))
		name := strings.Split(node.Name, ".")[0]
		instanceId := node.Spec.ProviderID[strings.LastIndex(node.Spec.ProviderID, "/")+1:]

		labels := node.GetLabels()

		nodegroup, ok := labels["eks.amazonaws.com/nodegroup"]
		if !ok {
			nodegroup = "-"
		}

		zone, ok := labels["topology.kubernetes.io/zone"]
		if !ok {
			zone = "unknown"
		}

		instanceType, ok := labels["node.kubernetes.io/instance-type"]
		if !ok {
			instanceType = "unknown"
		}

		table.AppendRow([]string{
			age.String(),
			name,
			instanceId,
			zone,
			nodegroup,
			instanceType,
		})
	}

	table.Print(writer)
	if len(p.nodes) > 0 {
		node := p.nodes[0]
		nodeSuffix := node.Name[strings.Index(node.Name, "."):]
		fmt.Printf("* Names end with %q\n", nodeSuffix)
	}

	return nil
}

func (p *NodePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.nodes)
}

func (p *NodePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.nodes)
}
