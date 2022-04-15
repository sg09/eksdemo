package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func ASGDescribeAutoScalingGroups(name string) ([]*autoscaling.Group, error) {
	sess := GetSession()
	svc := autoscaling.New(sess)

	autoScalingGroups := []*autoscaling.Group{}
	input := &autoscaling.DescribeAutoScalingGroupsInput{}
	pageNum := 0

	if name != "" {
		input.AutoScalingGroupNames = aws.StringSlice([]string{name})
	}

	err := svc.DescribeAutoScalingGroupsPages(input,
		func(page *autoscaling.DescribeAutoScalingGroupsOutput, lastPage bool) bool {
			pageNum++
			autoScalingGroups = append(autoScalingGroups, page.AutoScalingGroups...)
			return pageNum <= maxPages
		},
	)

	if err != nil {
		return nil, err
	}

	return autoScalingGroups, nil
}
