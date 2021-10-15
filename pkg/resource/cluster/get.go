package cluster

import (
	"eksdemo/pkg/aws"

	"github.com/aws/aws-sdk-go/service/eks"
)

func Get() ([]*eks.Cluster, error) {
	clusterNames, err := aws.EksListClusters()
	clusters := make([]*eks.Cluster, 0, len(clusterNames))

	if err != nil {
		return nil, err
	}

	for _, name := range clusterNames {
		result, err := aws.EksDescribeCluster(aws.StringValue(name))
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, result)
	}

	return clusters, nil
}
