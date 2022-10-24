package aws

import (
	"context"
	"time"

	cloudformationv2 "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	"github.com/aws/aws-sdk-go/aws"
)

type CloudformationClient struct {
	*cloudformationv2.Client
}

func NewCloudformationClient() *CloudformationClient {
	return &CloudformationClient{cloudformationv2.NewFromConfig(GetConfig())}
}

func (c *CloudformationClient) CreateStack(stackName, templateBody string, params map[string]string, caps []types.Capability) error {
	_, err := c.Client.CreateStack(context.Background(), &cloudformationv2.CreateStackInput{
		Capabilities: caps,
		Parameters:   toCloudformationParameters(params),
		StackName:    aws.String(stackName),
		TemplateBody: aws.String(templateBody),
	})
	if err != nil {
		return err
	}

	waiter := cloudformationv2.NewStackCreateCompleteWaiter(c.Client, func(o *cloudformationv2.StackCreateCompleteWaiterOptions) {
		o.APIOptions = append(o.APIOptions, WaiterLogger{}.AddLogger)
		o.MinDelay = 2 * time.Second
		o.MaxDelay = 5 * time.Second
	})

	return waiter.Wait(context.Background(),
		&cloudformationv2.DescribeStacksInput{StackName: aws.String(stackName)},
		2*time.Minute,
	)
}

func (c *CloudformationClient) DeleteStack(stackName string) error {
	_, err := c.Client.DeleteStack(context.Background(), &cloudformationv2.DeleteStackInput{
		StackName: aws.String(stackName),
	})

	return err
}

func (c *CloudformationClient) DescribeStacks(stackName string) ([]types.Stack, error) {
	input := cloudformationv2.DescribeStacksInput{}
	if stackName != "" {
		input.StackName = aws.String(stackName)
	}

	result, err := c.Client.DescribeStacks(context.Background(), &input)
	if err != nil {
		return nil, err
	}

	return result.Stacks, nil
}

func toCloudformationParameters(tags map[string]string) (params []types.Parameter) {
	for k, v := range tags {
		params = append(params, types.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		})
	}
	return
}
