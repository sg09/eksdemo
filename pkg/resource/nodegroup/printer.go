package nodegroup

import (
	"eksdemo/pkg/printer"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/hako/durafmt"
)

type NodegroupPrinter struct {
	Nodegroups []*eks.Nodegroup
}

func NewPrinter(Nodegroups []*eks.Nodegroup) *NodegroupPrinter {
	return &NodegroupPrinter{Nodegroups}
}

func (p *NodegroupPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "NodeGroup", "Version", "Type", "Nodes", "Instance(s)"})

	for _, n := range p.Nodegroups {
		age := durafmt.ParseShort(time.Since(*n.CreatedAt))

		table.AppendRow([]string{
			age.String(),
			*n.Status,
			*n.NodegroupName,
			*n.ReleaseVersion,
			*n.CapacityType,
			strconv.FormatInt(*n.ScalingConfig.DesiredSize, 10),
			strings.Join(aws.StringValueSlice(n.InstanceTypes), ","),
		})
	}

	table.Print(writer)

	return nil
}

func (p *NodegroupPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.Nodegroups)
}

func (p *NodegroupPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.Nodegroups)
}
