package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func ELBDeleteLoadBalancerV1(name string) error {
	sess := GetSession()
	svc := elb.New(sess)

	input := &elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(name),
	}
	_, err := svc.DeleteLoadBalancer(input)

	return err
}

func ELBDeleteLoadBalancerV2(arn string) error {
	sess := GetSession()
	svc := elbv2.New(sess)

	input := &elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: aws.String(arn),
	}
	_, err := svc.DeleteLoadBalancer(input)

	return err
}

func ELBDeleteTargetGroup(arn string) error {
	sess := GetSession()
	svc := elbv2.New(sess)

	input := &elbv2.DeleteTargetGroupInput{
		TargetGroupArn: aws.String(arn),
	}
	_, err := svc.DeleteTargetGroup(input)

	return err
}

func ELBDescribeListeners(loadBalancerArn string) ([]*elbv2.Listener, error) {
	sess := GetSession()
	svc := elbv2.New(sess)

	listeners := []*elbv2.Listener{}
	input := &elbv2.DescribeListenersInput{}
	pageNum := 0

	if loadBalancerArn != "" {
		input.LoadBalancerArn = aws.String(loadBalancerArn)
	}

	err := svc.DescribeListenersPages(input,
		func(page *elbv2.DescribeListenersOutput, lastPage bool) bool {
			pageNum++
			listeners = append(listeners, page.Listeners...)
			return pageNum <= maxPages
		},
	)

	return listeners, err
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

func ELBDescribeRules(listenerArn string, ruleArns []string) ([]*elbv2.Rule, error) {
	sess := GetSession()
	svc := elbv2.New(sess)

	input := &elbv2.DescribeRulesInput{}

	if listenerArn != "" {
		input.ListenerArn = aws.String(listenerArn)
	} else {
		input.RuleArns = aws.StringSlice(ruleArns)
	}

	result, err := svc.DescribeRules(input)

	return result.Rules, err
}

func ELBDescribeTargetGroups(name, loadBalancerArn string) ([]*elbv2.TargetGroup, error) {
	sess := GetSession()
	svc := elbv2.New(sess)

	targetGroups := []*elbv2.TargetGroup{}
	input := &elbv2.DescribeTargetGroupsInput{}
	pageNum := 0

	if name != "" {
		input.Names = aws.StringSlice([]string{name})
	}

	if loadBalancerArn != "" {
		input.LoadBalancerArn = aws.String(loadBalancerArn)
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

func ELBDescribeTargetHealth(arn, id string) ([]*elbv2.TargetHealthDescription, error) {
	sess := GetSession()
	svc := elbv2.New(sess)

	input := &elbv2.DescribeTargetHealthInput{
		TargetGroupArn: aws.String(arn),
	}

	if id != "" {
		input.Targets = []*elbv2.TargetDescription{
			{
				Id: aws.String(id),
			},
		}
	}

	res, err := svc.DescribeTargetHealth(input)

	return res.TargetHealthDescriptions, err
}
