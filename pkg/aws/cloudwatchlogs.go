package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type CloudwatchlogsClient struct {
	*cloudwatchlogs.Client
}

func NewCloudwatchlogsClient() *CloudwatchlogsClient {
	return &CloudwatchlogsClient{cloudwatchlogs.NewFromConfig(GetConfig())}
}

// Creates a log group with the specified name.
func (c *CloudwatchlogsClient) CreateLogGroup(name string) (*cloudwatchlogs.CreateLogGroupOutput, error) {
	return c.Client.CreateLogGroup(context.Background(), &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(name),
	})
}

// Deletes the specified log group and permanently deletes all the archived log events associated with the log group.
func (c *CloudwatchlogsClient) DeleteLogGroup(name string) error {
	_, err := c.Client.DeleteLogGroup(context.Background(), &cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: aws.String(name),
	})

	return err
}

// Lists the specified log groups. You can list all your log groups or filter the results by prefix.
// The results are ASCII-sorted by log group name.
func (c *CloudwatchlogsClient) DescribeLogGroups(namePrefix string) ([]types.LogGroup, error) {
	logGroups := []types.LogGroup{}
	pageNum := 0

	input := cloudwatchlogs.DescribeLogGroupsInput{}
	if namePrefix != "" {
		input.LogGroupNamePrefix = aws.String(namePrefix)
	}

	paginator := cloudwatchlogs.NewDescribeLogGroupsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		logGroups = append(logGroups, out.LogGroups...)
		pageNum++
	}

	return logGroups, nil
}

// Lists the log streams for the specified log group. You can list all the log streams or filter the results by prefix.
// You can also control how the results are ordered.
func (c *CloudwatchlogsClient) DescribeLogStreams(namePrefix, logGroupName string) ([]types.LogStream, error) {
	logStreams := []types.LogStream{}
	pageNum := 0

	input := cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(logGroupName),
	}

	if namePrefix != "" {
		input.LogStreamNamePrefix = aws.String(namePrefix)
	}

	paginator := cloudwatchlogs.NewDescribeLogStreamsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		logStreams = append(logStreams, out.LogStreams...)
		pageNum++
	}

	return logStreams, nil
}

// Lists log events from the specified log stream. You can list all of the log events or filter using a time range.
// By default, this operation returns as many log events as can fit in a response size of 1MB (up to 10,000 log events).
func (c *CloudwatchlogsClient) GetLogEvents(logStreamName, logGroupName string) ([]types.OutputLogEvent, error) {
	logEvents := []types.OutputLogEvent{}
	pageNum := 0

	input := cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
		StartFromHead: aws.Bool(true),
	}

	paginator := cloudwatchlogs.NewGetLogEventsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		logEvents = append(logEvents, out.Events...)
		pageNum++
	}

	return logEvents, nil
}
