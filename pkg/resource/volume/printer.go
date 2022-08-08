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

type VolumePrinter struct {
	volumes []*ec2.Volume
}

func NewPrinter(volumes []*ec2.Volume) *VolumePrinter {
	return &VolumePrinter{volumes}
}

func (p *VolumePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "State", "Id", "Name", "Type", "GiB", "AZ"})

	for _, v := range p.volumes {
		age := durafmt.ParseShort(time.Since(aws.TimeValue(v.CreateTime)))

		table.AppendRow([]string{
			age.String(),
			aws.StringValue(v.State),
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

func (p *VolumePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.volumes)
}

func (p *VolumePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.volumes)
}

func (p *VolumePrinter) getVolumeName(instance *ec2.Volume) string {
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
