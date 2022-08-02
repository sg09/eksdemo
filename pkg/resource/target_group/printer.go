package target_group

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"io"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/service/elbv2"
)

type TargetGroupPrinter struct {
	targetGroups []*elbv2.TargetGroup
}

func NewPrinter(targetGroups []*elbv2.TargetGroup) *TargetGroupPrinter {
	return &TargetGroupPrinter{targetGroups}
}

func (p *TargetGroupPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Name", "Port", "Proto", "Target", "Load Balancer"})

	for _, tg := range p.targetGroups {
		table.AppendRow([]string{
			aws.StringValue(tg.TargetGroupName),
			strconv.FormatInt(aws.Int64Value(tg.Port), 10),
			aws.StringValue(tg.Protocol),
			aws.StringValue(tg.TargetType),
			getLoadBalancer(tg),
		})
	}

	table.Print(writer)

	return nil
}

func (p *TargetGroupPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.targetGroups)
}

func (p *TargetGroupPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.targetGroups)
}

func getLoadBalancer(tg *elbv2.TargetGroup) string {
	lbArns := tg.LoadBalancerArns
	if len(lbArns) == 0 {
		return "None associated"
	} else if len(lbArns) > 1 {
		return strconv.Itoa(len(lbArns)) + " LBs associated"
	}

	splitArn := strings.Split(aws.StringValue(lbArns[0]), "/")
	if len(splitArn) < 2 {
		return "Error parsing LoadBalancer ARN"
	}
	return splitArn[2]
}
