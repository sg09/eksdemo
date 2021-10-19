package util

import (
	"eksdemo/pkg/aws"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

// kubernetes.io/cluster/<clusterName>
const K8stag = `kubernetes.io/cluster/%s`

func getSubnets(clusterName string) (string, error) {
	stackName := "eksctl-" + clusterName + "-cluster"

	stack, err := aws.CloudFormationDescribeStack(stackName)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case "ValidationError":
				return "", fmt.Errorf("cloudformation stack %q not found, is this an eksctl cluster?", stackName)
			default:
				return "", err
			}
		}
		return "", err
	}

	subnets := ""
	for _, o := range stack.Outputs {
		if aws.StringValue(o.OutputKey) == "SubnetsPrivate" {
			subnets = aws.StringValue(o.OutputValue)
		}
	}

	if subnets == "" {
		return "", fmt.Errorf("no private subnets found in cloudformation stack %q", stackName)
	}

	return subnets, nil
}

func CheckSubnets(clusterName string) error {
	subnetsString, err := getSubnets(clusterName)
	if err != nil {
		return err
	}

	subnets := strings.Split(subnetsString, ",")
	tag := fmt.Sprintf(K8stag, clusterName)
	tagsFilter := []string{tag}

	result, err := aws.EC2DescribeTags(subnets, tagsFilter)
	if err != nil {
		return err
	}

	if len(result.Tags) == 0 {
		return fmt.Errorf("required tag %q not found on any of the following private subnets:\n%s", tag, strings.Join(subnets, "\n"))
	}

	return nil
}

func TagSubnets(clusterName string) error {
	subnets, err := getSubnets(clusterName)
	if err != nil {
		return err
	}

	tags := map[string]string{
		fmt.Sprintf(K8stag, clusterName): "",
	}

	fmt.Println("Tagging subnets: " + subnets)
	fmt.Printf("With: %s\n", tags)

	return aws.EC2CreateTags(strings.Split(subnets, ","), tags)
}
