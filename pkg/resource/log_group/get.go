package log_group

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type Getter struct {
	resource.EmptyInit
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	logGroups, err := aws.CloudWatchLogsDescribeLogGroups(name)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(logGroups))
}

func (g *Getter) GetLogGroupByName(name string) (*cloudwatchlogs.LogGroup, error) {
	logGroups, err := aws.CloudWatchLogsDescribeLogGroups(name)
	if err != nil {
		// Return all errors except NotFound
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case cloudwatchlogs.ErrCodeResourceNotFoundException:
				break
			default:
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if len(logGroups) > 1 {
		return nil, fmt.Errorf("multiple log groups found with name %q", name)
	}

	if len(logGroups) == 0 {
		return nil, fmt.Errorf("log group %q not found", name)
	}

	return logGroups[0], nil
}
