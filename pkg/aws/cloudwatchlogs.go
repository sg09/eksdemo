package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func CloudWatchLogsDeleteLogGroup(name string) error {
	sess := GetSession()
	svc := cloudwatchlogs.New(sess)

	_, err := svc.DeleteLogGroup(&cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: aws.String(name),
	})

	return err
}

func CloudWatchLogsDescribeLogGroups(namePrefix string) ([]*cloudwatchlogs.LogGroup, error) {
	sess := GetSession()
	svc := cloudwatchlogs.New(sess)

	logGroups := []*cloudwatchlogs.LogGroup{}
	pageNum := 0

	input := &cloudwatchlogs.DescribeLogGroupsInput{}
	if namePrefix != "" {
		input.LogGroupNamePrefix = aws.String(namePrefix)
	}

	err := svc.DescribeLogGroupsPages(input,
		func(page *cloudwatchlogs.DescribeLogGroupsOutput, lastPage bool) bool {
			pageNum++
			logGroups = append(logGroups, page.LogGroups...)
			return pageNum <= maxPages
		})

	if err != nil {
		return nil, err
	}
	return logGroups, nil
}

func CloudWatchLogsDescribeLogStreams(namePrefix, logGroupName string) ([]*cloudwatchlogs.LogStream, error) {
	sess := GetSession()
	svc := cloudwatchlogs.New(sess)

	logStreams := []*cloudwatchlogs.LogStream{}
	pageNum := 0

	input := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(logGroupName),
	}

	if namePrefix != "" {
		input.LogStreamNamePrefix = aws.String(namePrefix)
	}

	err := svc.DescribeLogStreamsPages(input,
		func(page *cloudwatchlogs.DescribeLogStreamsOutput, lastPage bool) bool {
			pageNum++
			logStreams = append(logStreams, page.LogStreams...)
			return pageNum <= maxPages
		})

	if err != nil {
		return nil, err
	}
	return logStreams, nil
}
