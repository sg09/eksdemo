package log_stream

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/log_group"
	"fmt"
	"os"
)

type Getter struct {
	resource.EmptyInit
	logGroupGetter log_group.Getter
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	lsOptions, ok := options.(*LogStreamOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to LogStreamOptions")
	}

	logGroup, err := g.logGroupGetter.GetLogGroupByName(lsOptions.LogGroupName)
	if err != nil {
		return err
	}

	logStreams, err := aws.CloudWatchLogsDescribeLogStreams(name, aws.StringValue(logGroup.LogGroupName))
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(logStreams))
}
