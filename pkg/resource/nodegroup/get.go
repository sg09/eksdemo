package nodegroup

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"os"

	"github.com/aws/aws-sdk-go/service/eks"
)

type Getter struct {
	resource.EmptyInit
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var nodeGroup *eks.Nodegroup
	var nodeGroups []*eks.Nodegroup
	var err error

	clusterName := options.Common().ClusterName

	if name != "" {
		nodeGroup, err = g.GetNodeGroupByName(name, clusterName)
		nodeGroups = []*eks.Nodegroup{nodeGroup}
	} else {
		nodeGroups, err = g.GetAllNodeGroups(clusterName)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(nodeGroups))
}

func (g *Getter) GetAllNodeGroups(clusterName string) ([]*eks.Nodegroup, error) {
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

func (g *Getter) GetNodeGroupByName(name, clusterName string) (*eks.Nodegroup, error) {
	nodeGroup, err := aws.EksDescribeNodegroup(clusterName, name)
	if err != nil {
		return nil, err
	}

	return nodeGroup, nil
}
