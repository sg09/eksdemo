package util

import (
	"eksdemo/pkg/aws"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

func TagSubnets(clusterName string) error {
	stackName := "eksctl-" + clusterName + "-cluster"

	stack, err := aws.CloudFormationDescribeStack(stackName)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case "ValidationError":
				return fmt.Errorf("cloudformation stack %q not found, is this an eksctl cluster?", stackName)
			default:
				return err
			}
		}
		return err
	}

	subnets := ""
	for _, o := range stack.Outputs {
		if aws.StringValue(o.OutputKey) == "SubnetsPrivate" {
			subnets = aws.StringValue(o.OutputValue)
		}
	}

	if subnets == "" {
		return fmt.Errorf("no private subnets found in cloudformation stack %q", stackName)
	}

	tags := map[string]string{
		"kubernetes.io/cluster/" + clusterName: "",
	}

	fmt.Println("Tagging subnets: " + subnets)
	fmt.Printf("With: %s\n", tags)

	return aws.EC2CreateTags(strings.Split(subnets, ","), tags)
}
