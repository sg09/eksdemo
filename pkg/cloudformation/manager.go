package cloudformation

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type ResourceManager struct {
	Resource string

	Capabilities []aws.Capability
	DryRun       bool
	Parameters   map[string]string
	Template     template.Template
}

// eksdemo-<clusterName>-<resourceName>
const stackNameTemplate = "eksdemo-%s-%s"

func (m *ResourceManager) Create(options resource.Options) error {
	cfTemplate, err := m.Template.Render(options)
	if err != nil {
		return err
	}

	stackName := fmt.Sprintf(stackNameTemplate, options.Common().ClusterName, options.Common().Name)

	if m.DryRun {
		fmt.Println("\nCloudFormation Resource Manager Dry Run:")
		fmt.Printf("Stack name %q template:\n", stackName)
		fmt.Println(cfTemplate)
		return nil
	}

	fmt.Printf("Creating Cloudformation stack %q (can take 30+ seconds)...", stackName)
	err = aws.CloudFormationCreateStack(stackName, cfTemplate, m.Parameters, m.Capabilities)
	fmt.Println()

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case cloudformation.ErrCodeAlreadyExistsException:
				fmt.Printf("Cloudformation stack %q already exists\n", stackName)
				return nil
			default:
				return err
			}
		}
	}
	fmt.Printf("Cloudformation stack %q created\n", stackName)

	return nil
}

func (m *ResourceManager) Delete(options resource.Options) error {
	stackName := fmt.Sprintf(stackNameTemplate, options.Common().ClusterName, options.Common().Name)

	_, err := aws.CloudFormationDescribeStacks(stackName)
	if err != nil {
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				switch awsErr.Code() {
				case "ValidationError":
					fmt.Printf("Cloudformation stack %q doesn't exist, skipping...\n", stackName)
					return nil
				}
				return err
			}
		}
	}

	fmt.Printf("Deleting Cloudformation stack %q\n", stackName)

	return aws.CloudFormationDeleteStack(stackName)
}

func (m *ResourceManager) SetDryRun() {
	m.DryRun = true
}

func (m *ResourceManager) Update(options resource.Options) error {
	return fmt.Errorf("feature not supported")
}
