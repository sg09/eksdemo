package volume

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hako/durafmt"
)

const maxNameLength = 25

type EC2Printer struct {
	volumes []*ec2.Volume
}

func NewPrinter(volumes []*ec2.Volume) *EC2Printer {
	return &EC2Printer{volumes}
}

func (p *EC2Printer) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Id", "Name", "Type", "GiB", "AZ"})

	for _, v := range p.volumes {
		age := durafmt.ParseShort(time.Since(aws.TimeValue(v.CreateTime)))

		table.AppendRow([]string{
			age.String(),
			aws.StringValue(v.VolumeId),
			p.getVolumeName(v),
			aws.StringValue(v.VolumeType),
			strconv.FormatInt(aws.Int64Value(v.Size), 10),
			aws.StringValue(v.AvailabilityZone),
		})
	}
	table.Print(writer)

	return nil
}

func (p *EC2Printer) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.volumes)
}

func (p *EC2Printer) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.volumes)
}

func (p *EC2Printer) getVolumeName(instance *ec2.Volume) string {
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
