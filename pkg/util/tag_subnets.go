package util

import (
	"eksdemo/pkg/aws"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

// kubernetes.io/cluster/<clusterName>
const K8stag = `kubernetes.io/cluster/%s`

func GetPrivateSubnets(clusterName string) ([]string, error) {
	stackName := "eksctl-" + clusterName + "-cluster"

	stacks, err := aws.NewCloudformationClient().DescribeStacks(stackName)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case "ValidationError":
				return nil, fmt.Errorf("cloudformation stack %q not found, is this an eksctl cluster?", stackName)
			default:
				return nil, err
			}
		}
		return nil, err
	}

	subnets := ""
	for _, o := range stacks[0].Outputs {
		if aws.StringValue(o.OutputKey) == "SubnetsPrivate" {
			subnets = aws.StringValue(o.OutputValue)
			continue
		}
	}

	if subnets == "" {
		return nil, fmt.Errorf("no private subnets found in cloudformation stack %q", stackName)
	}

	return strings.Split(subnets, ","), nil
}

func CheckSubnets(clusterName string) error {
	subnets, err := GetPrivateSubnets(clusterName)
	if err != nil {
		return err
	}

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
	subnets, err := GetPrivateSubnets(clusterName)
	if err != nil {
		return err
	}

	tags := map[string]string{
		fmt.Sprintf(K8stag, clusterName): "",
	}

	fmt.Println("Tagging subnets: " + strings.Join(subnets, ","))
	fmt.Printf("With: %s\n", tags)

	return aws.EC2CreateTags(subnets, tags)
}
