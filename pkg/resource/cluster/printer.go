package cluster

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/hako/durafmt"

	"github.com/aws/aws-sdk-go/service/eks"
)

type ClusterPrinter struct {
	clusters   []*eks.Cluster
	clusterURL string
}

func NewPrinter(clusters []*eks.Cluster, clusterURL string) *ClusterPrinter {
	return &ClusterPrinter{clusters, clusterURL}
}

func (p *ClusterPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Cluster", "Version", "Platform", "Endpoint", "Logging"})
	currentContext := false

	for _, cluster := range p.clusters {
		var endpoint string

		vpcConf := cluster.ResourcesVpcConfig
		if vpcConf == nil {
			endpoint = "-"
		} else if *vpcConf.EndpointPublicAccess && !*vpcConf.EndpointPrivateAccess {
			endpoint = "Public"
		} else if *vpcConf.EndpointPublicAccess && *vpcConf.EndpointPrivateAccess {
			endpoint = "Public/Private"
		} else {
			endpoint = "Private"
		}

		age := durafmt.ParseShort(time.Since(*cluster.CreatedAt))
		name := *cluster.Name

		if cluster.Endpoint != nil {
			if *cluster.Endpoint == p.clusterURL {
				currentContext = true
				name = "*" + name
			}
		}

		table.AppendRow([]string{
			age.String(),
			aws.StringValue(cluster.Status),
			name,
			aws.StringValue(cluster.Version),
			aws.StringValue(cluster.PlatformVersion),
			endpoint,
			strconv.FormatBool(*cluster.Logging.ClusterLogging[0].Enabled),
		})
	}

	table.Print(writer)
	if currentContext {
		fmt.Println("* Indicates current context")
	}

	return nil
}

func (p *ClusterPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.clusters)
}

func (p *ClusterPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.clusters)
}
