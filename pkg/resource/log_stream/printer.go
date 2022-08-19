package log_stream

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/hako/durafmt"
)

type LogStreamPrinter struct {
	logStreams []*cloudwatchlogs.LogStream
}

func NewPrinter(logStreams []*cloudwatchlogs.LogStream) *LogStreamPrinter {
	return &LogStreamPrinter{logStreams}
}

func (p *LogStreamPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Name", "Last Event"})

	for _, ls := range p.logStreams {
		age := durafmt.ParseShort(time.Since(time.Unix(aws.Int64Value(ls.CreationTime)/1000, 0)))
		lastEvent := durafmt.ParseShort(time.Since(time.Unix(aws.Int64Value(ls.LastEventTimestamp)/1000, 0)))

		table.AppendRow([]string{
			age.String(),
			aws.StringValue(ls.LogStreamName),
			lastEvent.String(),
		})
	}
	table.Print(writer)

	return nil
}

func (p *LogStreamPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.logStreams)
}

func (p *LogStreamPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.logStreams)
}
