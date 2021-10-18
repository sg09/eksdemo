package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func EksDescribeCluster(clusterName string) (*eks.Cluster, error) {
	sess := GetSession()
	svc := eks.New(sess)

	result, err := svc.DescribeCluster(&eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	})

	if err != nil {
		return nil, FormatError(err)
	}

	return result.Cluster, nil
}

func EksDescribeNodegroup(clusterName, nodegroupName string) (*eks.Nodegroup, error) {
	sess := GetSession()
	svc := eks.New(sess)

	result, err := svc.DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   aws.String(clusterName),
		NodegroupName: aws.String(nodegroupName),
	})

	if err != nil {
		return nil, FormatError(err)
	}

	return result.Nodegroup, nil
}

func EksListClusters() ([]*string, error) {
	sess := GetSession()
	svc := eks.New(sess)

	clusters := []*string{}
	pageNum := 0

	err := svc.ListClustersPages(&eks.ListClustersInput{},
		func(page *eks.ListClustersOutput, lastPage bool) bool {
			pageNum++
			clusters = append(clusters, page.Clusters...)
			return pageNum <= 3
		},
	)

	if err != nil {
		return nil, FormatError(err)
	}

	return clusters, nil
}

func EksListNodegroups(clusterName string) ([]*string, error) {
	sess := GetSession()
	svc := eks.New(sess)

	nodegroups, err := svc.ListNodegroups(&eks.ListNodegroupsInput{
		ClusterName: aws.String(clusterName),
	})

	if err != nil {
		return nil, FormatError(err)
	}

	return nodegroups.Nodegroups, nil
}

func EksOptimizedAmi(eksVersion string) (string, error) {
	sess := GetSession()
	svc := ssm.New(sess)

	input := ssm.GetParameterInput{
		Name: aws.String("/aws/service/eks/optimized-ami/" + eksVersion + "/amazon-linux-2/recommended/image_id"),
	}

	result, err := svc.GetParameter(&input)
	if err != nil {
		return "", fmt.Errorf("ssm failed to lookup EKS Optimized AMI: %s", err)
	}

	return aws.StringValue(result.Parameter.Value), nil
}
