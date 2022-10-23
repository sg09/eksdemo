package cluster

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/kubernetes"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"os"

	"github.com/aws/aws-sdk-go/service/eks"
)

type Getter struct {
	resource.EmptyInit
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var clusters []*eks.Cluster
	var err error

	if name != "" {
		clusters, err = g.GetClusterByName(name)
	} else {
		clusters, err = g.GetAllClusters()
	}

	if err != nil {
		return err
	}

	currentClusterUrl := kubernetes.ClusterURLForCurrentContext()

	return output.Print(os.Stdout, NewPrinter(clusters, currentClusterUrl))
}

func (g *Getter) GetAllClusters() ([]*eks.Cluster, error) {
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

func (g *Getter) GetClusterByName(name string) ([]*eks.Cluster, error) {
	cluster, err := aws.EksDescribeCluster(name)
	if err != nil {
		return nil, err
	}

	return []*eks.Cluster{cluster}, nil
}
