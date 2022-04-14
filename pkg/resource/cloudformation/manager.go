package cloudformation

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

type Manager struct {
	DryRun bool
}

func (m *Manager) Delete(options resource.Options) error {
	stackName := options.Common().Name

	_, err := aws.CloudFormationDescribeStacks(stackName)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case "ValidationError":
				return fmt.Errorf(awsErr.Message())
			}
			return err
		}
	}

	fmt.Printf("Deleting Cloudformation stack %q\n", stackName)

	return aws.CloudFormationDeleteStack(stackName)
}

func (m *Manager) Create(options resource.Options) error {
	return fmt.Errorf("feature not yet implemented")
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options) error {
	return fmt.Errorf("feature not supported")
}
