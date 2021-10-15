package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func CloudFormationDescribeStack(stackName string) (*cloudformation.Stack, error) {
	sess := GetSession()
	svc := cloudformation.New(sess)

	result, err := svc.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	})

	if err != nil {
		return nil, FormatError(err)
	}

	return result.Stacks[0], nil
}
