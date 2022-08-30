package ec2_instance

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hako/durafmt"
)

const maxNameLength = 30

type EC2Printer struct {
	reservations []*ec2.Reservation
}

func NewPrinter(reservations []*ec2.Reservation) *EC2Printer {
	return &EC2Printer{reservations}
}

func (p *EC2Printer) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "State", "Id", "Name", "Type", "Zone"})
	spot := 0

	for _, res := range p.reservations {
		for _, i := range res.Instances {
			age := durafmt.ParseShort(time.Since(aws.TimeValue(i.LaunchTime)))

			instanceType := aws.StringValue(i.InstanceType)
			if aws.StringValue(i.InstanceLifecycle) == "spot" {
				instanceType = "*" + instanceType
				spot++
			}

			table.AppendRow([]string{
				age.String(),
				aws.StringValue(i.State.Name),
				aws.StringValue(i.InstanceId),
				p.getInstanceName(i),
				instanceType,
				aws.StringValue(i.Placement.AvailabilityZone),
			})
		}
	}

	table.Print(writer)
	if spot > 0 {
		fmt.Println("* Indicates Spot Instance")
	}

	return nil
}

func (p *EC2Printer) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.reservations)
}

func (p *EC2Printer) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.reservations)
}

func (p *EC2Printer) getInstanceName(instance *ec2.Instance) string {
	name := ""
	for _, tag := range instance.Tags {
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
