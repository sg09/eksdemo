package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type CloudWatchLogsClient struct {
	*cloudwatchlogs.Client
}

func NewCloudWatchLogsClient() *CloudWatchLogsClient {
	return &CloudWatchLogsClient{cloudwatchlogs.NewFromConfig(GetConfig())}
}

func (c *CloudWatchLogsClient) DeleteLogGroup(name string) error {
	_, err := c.Client.DeleteLogGroup(context.Background(), &cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: aws.String(name),
	})

	return err
}

func (c *CloudWatchLogsClient) DescribeLogGroups(namePrefix string) ([]types.LogGroup, error) {
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

func (c *CloudWatchLogsClient) DescribeLogStreams(namePrefix, logGroupName string) ([]types.LogStream, error) {
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

func (c *CloudWatchLogsClient) GetLogEvents(logStreamName, logGroupName string) ([]types.OutputLogEvent, error) {
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
