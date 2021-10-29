package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type Capability string

const (
	CapabilityCapabilityIam        Capability = "CAPABILITY_IAM"
	CapabilityCapabilityNamedIam   Capability = "CAPABILITY_NAMED_IAM"
	CapabilityCapabilityAutoExpand Capability = "CAPABILITY_AUTO_EXPAND"
)

func CloudFormationCreateStack(stackName, templateBody string, parameters map[string]string, capabilities []Capability) error {
	sess := GetSession()
	svc := cloudformation.New(sess)

	_, err := svc.CreateStack(&cloudformation.CreateStackInput{
		Capabilities: createCloudFormationCapabilities(capabilities),
		Parameters:   createCloudFormationParameters(parameters),
		StackName:    aws.String(stackName),
		TemplateBody: aws.String(templateBody),
	})
	if err != nil {
		return err
	}

	err = svc.WaitUntilStackCreateComplete(&cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return err
	}

	return nil
}

func CloudFormationDeleteStack(stackName string) error {
	sess := GetSession()
	svc := cloudformation.New(sess)

	_, err := CloudFormationDescribeStack(stackName)
	if err != nil {
		return err
	}

	_, err = svc.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return err
	}

	return nil
}

func CloudFormationDescribeStack(stackName string) (*cloudformation.Stack, error) {
	sess := GetSession()
	svc := cloudformation.New(sess)

	result, err := svc.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	})

	if err != nil {
		return nil, err
	}

	return result.Stacks[0], nil
}

func createCloudFormationParameters(tags map[string]string) (cfParams []*cloudformation.Parameter) {
	for k, v := range tags {
		cfParams = append(cfParams, &cloudformation.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		})
	}
	return
}

func createCloudFormationCapabilities(c []Capability) []*string {
	s := make([]*string, len(c))
	for i, v := range c {
		s[i] = aws.String(string(v))
	}
	return s
}
