package log_event

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/log_group"
	"fmt"
	"os"
)

type Getter struct {
	logGroupGetter log_group.Getter
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	lsOptions, ok := options.(*LogEventOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to LogStreamOptions")
	}

	logGroup, err := g.logGroupGetter.GetLogGroupByName(lsOptions.LogGroupName)
	if err != nil {
		return err
	}

	logEvents, err := aws.CloudWatchLogsGetLogEvents(name, aws.StringValue(logGroup.LogGroupName))
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(logEvents, lsOptions.Timestamp))
}
