package log_event

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/hako/durafmt"
)

type LogEventPrinter struct {
	logEvents []*cloudwatchlogs.OutputLogEvent
	timestamp bool
}

func NewPrinter(logEvents []*cloudwatchlogs.OutputLogEvent, timestamp bool) *LogEventPrinter {
	return &LogEventPrinter{logEvents, timestamp}
}

func (p *LogEventPrinter) PrintTable(writer io.Writer) error {
	for _, le := range p.logEvents {
		if p.timestamp {
			age := durafmt.ParseShort(time.Since(time.Unix(aws.Int64Value(le.Timestamp)/1000, 0)))
			fmt.Printf(strings.ReplaceAll(age.InternationalString(), " ", "") + ": ")
		}
		// strings.Join and strings.Fields combination removes extra spaces
		fmt.Println(strings.Join(strings.Fields(aws.StringValue(le.Message)), " "))
	}

	return nil
}

func (p *LogEventPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.logEvents)
}

func (p *LogEventPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.logEvents)
}
