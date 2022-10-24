package auto_scaling_group

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"strconv"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/hako/durafmt"
)

type AutoScalingGroupPrinter struct {
	autoScalingGroups []types.AutoScalingGroup
}

func NewPrinter(zones []types.AutoScalingGroup) *AutoScalingGroupPrinter {
	return &AutoScalingGroupPrinter{zones}
}

func (p *AutoScalingGroupPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Name", "Instances", "Desired", "Min", "Max"})

	for _, asg := range p.autoScalingGroups {
		age := durafmt.ParseShort(time.Since(aws.TimeValue(asg.CreatedTime)))

		table.AppendRow([]string{
			age.String(),
			awssdk.ToString(asg.AutoScalingGroupName),
			strconv.Itoa(len(asg.Instances)),
			strconv.Itoa(int(awssdk.ToInt32(asg.DesiredCapacity))),
			strconv.Itoa(int(awssdk.ToInt32(asg.MinSize))),
			strconv.Itoa(int(awssdk.ToInt32(asg.MaxSize))),
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
