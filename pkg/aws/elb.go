package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func ELBDeleteTargetGroup(arn string) error {
	sess := GetSession()
	svc := elbv2.New(sess)

	input := &elbv2.DeleteTargetGroupInput{
		TargetGroupArn: aws.String(arn),
	}
	_, err := svc.DeleteTargetGroup(input)

	return err
}

func ELBDescribeLoadBalancersv1(name string) ([]*elb.LoadBalancerDescription, error) {
	sess := GetSession()
	svc := elb.New(sess)

	elbs := []*elb.LoadBalancerDescription{}
	input := &elb.DescribeLoadBalancersInput{}
	pageNum := 0

	if name != "" {
		input.LoadBalancerNames = aws.StringSlice([]string{name})
	}

	err := svc.DescribeLoadBalancersPages(input,
		func(page *elb.DescribeLoadBalancersOutput, lastPage bool) bool {
			pageNum++
			elbs = append(elbs, page.LoadBalancerDescriptions...)
			return pageNum <= maxPages
		},
	)

	return elbs, err
}

func ELBDescribeLoadBalancersv2(name string) ([]*elbv2.LoadBalancer, error) {
	sess := GetSession()
	svc := elbv2.New(sess)

	elbs := []*elbv2.LoadBalancer{}
	input := &elbv2.DescribeLoadBalancersInput{}
	pageNum := 0

	if name != "" {
		input.Names = aws.StringSlice([]string{name})
	}

	err := svc.DescribeLoadBalancersPages(input,
		func(page *elbv2.DescribeLoadBalancersOutput, lastPage bool) bool {
			pageNum++
			elbs = append(elbs, page.LoadBalancers...)
			return pageNum <= maxPages
		},
	)

	return elbs, err
}

func ELBDescribeTargetGroups(name string) ([]*elbv2.TargetGroup, error) {
	sess := GetSession()
	svc := elbv2.New(sess)

	targetGroups := []*elbv2.TargetGroup{}
	input := &elbv2.DescribeTargetGroupsInput{}
	pageNum := 0

	if name != "" {
		input.Names = aws.StringSlice([]string{name})
	}

	err := svc.DescribeTargetGroupsPages(input,
		func(page *elbv2.DescribeTargetGroupsOutput, lastPage bool) bool {
			pageNum++
			targetGroups = append(targetGroups, page.TargetGroups...)
			return pageNum <= maxPages
		},
	)

	return targetGroups, err
}
