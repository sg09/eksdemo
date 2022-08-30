package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func SSMDescribeInstanceInformation(instanceId string) ([]*ssm.InstanceInformation, error) {
	sess := GetSession()
	svc := ssm.New(sess)

	filters := []*ssm.InstanceInformationStringFilter{}
	instances := []*ssm.InstanceInformation{}
	pageNum := 0

	if instanceId != "" {
		filters = append(filters, &ssm.InstanceInformationStringFilter{
			Key:    aws.String("InstanceIds"),
			Values: aws.StringSlice([]string{instanceId}),
		})
	}

	input := &ssm.DescribeInstanceInformationInput{}
	if len(filters) > 0 {
		input.Filters = filters
	}

	err := svc.DescribeInstanceInformationPages(input,
		func(page *ssm.DescribeInstanceInformationOutput, lastPage bool) bool {
			pageNum++
			instances = append(instances, page.InstanceInformationList...)
			return pageNum <= maxPages
		},
	)

	if err != nil {
		return nil, err
	}

	return instances, nil
}
