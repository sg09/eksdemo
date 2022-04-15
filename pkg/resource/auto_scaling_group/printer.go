package auto_scaling_group

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/hako/durafmt"
)

type AutoScalingGroupPrinter struct {
	autoScalingGroups []*autoscaling.Group
}

func NewPrinter(zones []*autoscaling.Group) *AutoScalingGroupPrinter {
	return &AutoScalingGroupPrinter{zones}
}

func (p *AutoScalingGroupPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Name", "Instances", "Desired", "Min", "Max"})

	for _, asg := range p.autoScalingGroups {
		age := durafmt.ParseShort(time.Since(aws.TimeValue(asg.CreatedTime)))

		table.AppendRow([]string{
			age.String(),
			aws.StringValue(asg.AutoScalingGroupName),
			strconv.Itoa(len(asg.Instances)),
			strconv.FormatInt(aws.Int64Value(asg.DesiredCapacity), 10),
			strconv.FormatInt(aws.Int64Value(asg.MinSize), 10),
			strconv.FormatInt(aws.Int64Value(asg.MaxSize), 10),
		})
	}

	table.Print(writer)

	return nil
}

func (p *AutoScalingGroupPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.autoScalingGroups)
}

func (p *AutoScalingGroupPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.autoScalingGroups)
}
