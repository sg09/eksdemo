package nodegroup

import (
	"eksdemo/pkg/aws"

	"github.com/aws/aws-sdk-go/service/eks"
)

func Get(clusterName string) ([]*eks.Nodegroup, error) {
	nodeGroupNames, err := aws.EksListNodegroups(clusterName)
	nodeGroups := make([]*eks.Nodegroup, 0, len(nodeGroupNames))

	if err != nil {
		return nil, err
	}

	for _, name := range nodeGroupNames {
		result, err := aws.EksDescribeNodegroup(clusterName, *name)
		if err != nil {
			return nil, err
		}
		nodeGroups = append(nodeGroups, result)
	}

	return nodeGroups, nil
}
