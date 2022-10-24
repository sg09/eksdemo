package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func EksDescribeAddon(clusterName, addonName string) (*eks.Addon, error) {
	sess := GetSession()
	svc := eks.New(sess)

	result, err := svc.DescribeAddon(&eks.DescribeAddonInput{
		AddonName:   aws.String(addonName),
		ClusterName: aws.String(clusterName),
	})

	if err != nil {
		return nil, err
	}

	return result.Addon, nil
}

func EksDescribeAddonVersions(addonName, version string) ([]*eks.AddonInfo, error) {
	sess := GetSession()
	svc := eks.New(sess)

	addons := []*eks.AddonInfo{}
	pageNum := 0

	input := &eks.DescribeAddonVersionsInput{
		KubernetesVersion: aws.String(version),
	}

	if addonName != "" {
		input.AddonName = aws.String(addonName)
	}

	err := svc.DescribeAddonVersionsPages(input,
		func(page *eks.DescribeAddonVersionsOutput, lastPage bool) bool {
			pageNum++
			addons = append(addons, page.Addons...)
			return pageNum <= maxPages
		},
	)

	if err != nil {
		return nil, err
	}

	return addons, nil
}

func EksDescribeCluster(clusterName string) (*eks.Cluster, error) {
	sess := GetSession()
	svc := eks.New(sess)

	result, err := svc.DescribeCluster(&eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	})

	if err != nil {
		return nil, FormatErrorSDKv1(err)
	}

	return result.Cluster, nil
}

func EksDescribeFargateProfile(clusterName, profileName string) (*eks.FargateProfile, error) {
	sess := GetSession()
	svc := eks.New(sess)

	result, err := svc.DescribeFargateProfile(&eks.DescribeFargateProfileInput{
		ClusterName:        aws.String(clusterName),
		FargateProfileName: aws.String(profileName),
	})

	if err != nil {
		return nil, FormatErrorSDKv1(err)
	}

	return result.FargateProfile, nil
}

func EksDescribeNodegroup(clusterName, nodegroupName string) (*eks.Nodegroup, error) {
	sess := GetSession()
	svc := eks.New(sess)

	result, err := svc.DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   aws.String(clusterName),
		NodegroupName: aws.String(nodegroupName),
	})

	if err != nil {
		return nil, FormatErrorSDKv1(err)
	}

	return result.Nodegroup, nil
}

func EksListAddons(clusterName string) ([]*string, error) {
	sess := GetSession()
	svc := eks.New(sess)

	addons := []*string{}
	pageNum := 0

	err := svc.ListAddonsPages(&eks.ListAddonsInput{
		ClusterName: aws.String(clusterName),
	},
		func(page *eks.ListAddonsOutput, lastPage bool) bool {
			pageNum++
			addons = append(addons, page.Addons...)
			return pageNum <= maxPages
		},
	)

	if err != nil {
		return nil, err
	}

	return addons, nil
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
			return pageNum <= maxPages
		},
	)

	if err != nil {
		return nil, FormatErrorSDKv1(err)
	}

	return clusters, nil
}

func EksListFargateProfiles(clusterName string) ([]*string, error) {
	sess := GetSession()
	svc := eks.New(sess)

	profiles := []*string{}
	pageNum := 0

	err := svc.ListFargateProfilesPages(&eks.ListFargateProfilesInput{
		ClusterName: aws.String(clusterName),
	},
		func(page *eks.ListFargateProfilesOutput, lastPage bool) bool {
			pageNum++
			profiles = append(profiles, page.FargateProfileNames...)
			return pageNum <= maxPages
		},
	)

	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func EksListNodegroups(clusterName string) ([]*string, error) {
	sess := GetSession()
	svc := eks.New(sess)

	nodegroups, err := svc.ListNodegroups(&eks.ListNodegroupsInput{
		ClusterName: aws.String(clusterName),
	})

	if err != nil {
		return nil, FormatErrorSDKv1(err)
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

func EksUpdateNodegroupConfig(clusterName, nodegroupName string, desired, min, max int) error {
	sess := GetSession()
	svc := eks.New(sess)

	_, err := svc.UpdateNodegroupConfig(&eks.UpdateNodegroupConfigInput{
		ClusterName:   aws.String(clusterName),
		NodegroupName: aws.String(nodegroupName),
		ScalingConfig: &eks.NodegroupScalingConfig{
			DesiredSize: aws.Int64(int64(desired)),
			MinSize:     aws.Int64(int64(min)),
			MaxSize:     aws.Int64(int64(max)),
		},
	})

	return err
}
