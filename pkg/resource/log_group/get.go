package log_group

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type Getter struct {
	cloudwatchlogsClient *aws.CloudwatchlogsClient
}

func NewGetter(cloudwatchlogsClient *aws.CloudwatchlogsClient) *Getter {
	return &Getter{cloudwatchlogsClient}
}

func (g *Getter) Init() {
	if g.cloudwatchlogsClient == nil {
		g.cloudwatchlogsClient = aws.NewCloudwatchlogsClient()
	}
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	logGroups, err := g.cloudwatchlogsClient.DescribeLogGroups(name)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(logGroups))
}

func (g *Getter) GetLogGroupByName(name string) (types.LogGroup, error) {
	logGroups, err := g.cloudwatchlogsClient.DescribeLogGroups(name)
	if err != nil {
		return types.LogGroup{}, err
	}

	if len(logGroups) > 1 {
		return types.LogGroup{}, fmt.Errorf("multiple log groups found with name %q", name)
	}

	if len(logGroups) == 0 {
		return types.LogGroup{}, fmt.Errorf("log group %q not found", name)
	}

	return logGroups[0], nil
}
