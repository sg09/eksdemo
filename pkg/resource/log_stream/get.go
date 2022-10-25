package log_stream

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/log_group"
	"fmt"
	"os"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
)

type Getter struct {
	cloudwatchlogsClient *aws.CloudWatchLogsClient
	logGroupGetter       *log_group.Getter
}

func NewGetter(cloudwatchlogsClient *aws.CloudWatchLogsClient) *Getter {
	return &Getter{cloudwatchlogsClient, log_group.NewGetter(cloudwatchlogsClient)}
}

func (g *Getter) Init() {
	if g.cloudwatchlogsClient == nil {
		g.cloudwatchlogsClient = aws.NewCloudWatchLogsClient()
	}
	g.logGroupGetter = log_group.NewGetter(g.cloudwatchlogsClient)
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

	logStreams, err := g.cloudwatchlogsClient.DescribeLogStreams(name, awssdk.ToString(logGroup.LogGroupName))
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(logStreams))
}
