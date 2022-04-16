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
	table.SetHeader([]string{"Age", "Status", "Name", "Nodes", "Min", "Max", "Version", "Type", "Instance(s)"})

	for _, n := range p.Nodegroups {
		age := durafmt.ParseShort(time.Since(aws.TimeValue(n.CreatedAt)))

		table.AppendRow([]string{
			age.String(),
			aws.StringValue(n.Status),
			aws.StringValue(n.NodegroupName),
			strconv.FormatInt(aws.Int64Value(n.ScalingConfig.DesiredSize), 10),
			strconv.FormatInt(aws.Int64Value(n.ScalingConfig.MinSize), 10),
			strconv.FormatInt(aws.Int64Value(n.ScalingConfig.MaxSize), 10),
			aws.StringValue(n.ReleaseVersion),
			aws.StringValue(n.CapacityType),
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
