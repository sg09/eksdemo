package cloudformation

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type Getter struct {
	resource.EmptyInit
}

func (g *Getter) Get(stackName string, output printer.Output, options resource.Options) error {
	var err error
	var stacks []*cloudformation.Stack
	clusterName := options.Common().ClusterName

	if clusterName != "" {
		stacks, err = g.GetStacksByCluster(clusterName, stackName)
	} else {
		stacks, err = g.GetStacks(stackName)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(stacks))
}

func (g *Getter) GetStacks(stackName string) ([]*cloudformation.Stack, error) {
	return aws.CloudFormationDescribeStacks(stackName)
}

func (g *Getter) GetStacksByCluster(clusterName, stackName string) ([]*cloudformation.Stack, error) {
	stacks, err := aws.CloudFormationDescribeStacks(stackName)

	if err != nil || clusterName == "" {
		return stacks, err
	}

	clusterStacks := make([]*cloudformation.Stack, 0, len(stacks))

	for _, stack := range stacks {
		name := aws.StringValue(stack.StackName)
		if strings.Contains(name, "eksdemo-"+clusterName+"-") || strings.Contains(name, "eksctl-"+clusterName+"-") {
			clusterStacks = append(clusterStacks, stack)
		}
	}

	return clusterStacks, nil
}
