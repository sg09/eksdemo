package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

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
