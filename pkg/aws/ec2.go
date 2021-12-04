package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func EC2CreateTags(resources []string, tags map[string]string) error {
	sess := GetSession()
	svc := ec2.New(sess)

	_, err := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: aws.StringSlice(resources),
		Tags:      createEC2Tags(tags),
	})

	if err != nil {
		return FormatError(err)
	}
	return nil
}

func EC2DescribeTags(resources, tagsFilter []string) (*ec2.DescribeTagsOutput, error) {
	sess := GetSession()
	svc := ec2.New(sess)

	filters := []*ec2.Filter{
		{
			Name:   aws.String("resource-id"),
			Values: aws.StringSlice(resources),
		},
	}

	if len(tagsFilter) > 0 {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("key"),
			Values: aws.StringSlice(tagsFilter),
		})
	}

	result, err := svc.DescribeTags(&ec2.DescribeTagsInput{
		Filters: filters,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func EC2DescribeSubnets(name, vpcId string) ([]*ec2.Subnet, error) {
	sess := GetSession()
	svc := ec2.New(sess)

	filters := []*ec2.Filter{}
	subnets := []*ec2.Subnet{}
	pageNum := 0

	if name != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("tag:Name"),
			Values: aws.StringSlice([]string{name}),
		})
	}

	if vpcId != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("vpc-id"),
			Values: aws.StringSlice([]string{vpcId}),
		})
	}

	input := &ec2.DescribeSubnetsInput{}

	if len(filters) > 0 {
		input.Filters = filters
	}

	err := svc.DescribeSubnetsPages(input,
		func(page *ec2.DescribeSubnetsOutput, lastPage bool) bool {
			pageNum++
			subnets = append(subnets, page.Subnets...)
			return pageNum <= maxPages
		},
	)

	if err != nil {
		return nil, err
	}

	return subnets, nil
}

func EC2DescribeVpcs(name, vpcId string) ([]*ec2.Vpc, error) {
	sess := GetSession()
	svc := ec2.New(sess)

	filters := []*ec2.Filter{}
	vpcs := []*ec2.Vpc{}
	pageNum := 0

	if name != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("tag:Name"),
			Values: aws.StringSlice([]string{name}),
		})
	}

	if vpcId != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("vpc-id"),
			Values: aws.StringSlice([]string{vpcId}),
		})
	}

	input := &ec2.DescribeVpcsInput{}

	if len(filters) > 0 {
		input.Filters = filters
	}

	err := svc.DescribeVpcsPages(input,
		func(page *ec2.DescribeVpcsOutput, lastPage bool) bool {
			pageNum++
			vpcs = append(vpcs, page.Vpcs...)
			return pageNum <= maxPages
		},
	)

	if err != nil {
		return nil, err
	}

	return vpcs, nil
}

func createEC2Tags(tags map[string]string) (ec2tags []*ec2.Tag) {
	for k, v := range tags {
		ec2tags = append(ec2tags, &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return
}
