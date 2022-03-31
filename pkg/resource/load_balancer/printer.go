package load_balancer

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/hako/durafmt"
)

type LoadBalancerPrinter struct {
	*LoadBalancers
}

func NewPrinter(loadBalancers *LoadBalancers) *LoadBalancerPrinter {
	return &LoadBalancerPrinter{loadBalancers}
}

func (p *LoadBalancerPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "State", "Name", "Type", "Stack", "AZs", "SGs"})

	for _, elb := range p.V1 {
		age := durafmt.ParseShort(time.Since(aws.TimeValue(elb.CreatedTime)))

		name := aws.StringValue(elb.LoadBalancerName)
		if aws.StringValue(elb.Scheme) == "internal" {
			name = "*" + name
		}

		table.AppendRow([]string{
			age.String(),
			"-",
			name,
			"CLB",
			"ipv4",
			strconv.Itoa(len(elb.AvailabilityZones)),
			strconv.Itoa(len(elb.SecurityGroups)),
		})
	}

	for _, elb := range p.V2 {
		age := durafmt.ParseShort(time.Since(aws.TimeValue(elb.CreatedTime)))

		elbType := "unknown"
		switch aws.StringValue(elb.Type) {
		case "application":
			elbType = "ALB"
		case "network":
			elbType = "NLB"
		case "gateway":
			elbType = "GWLB"
		}

		table.AppendRow([]string{
			age.String(),
			aws.StringValue(elb.State.Code),
			aws.StringValue(elb.LoadBalancerName),
			elbType,
			aws.StringValue(elb.IpAddressType),
			strconv.Itoa(len(elb.AvailabilityZones)),
			strconv.Itoa(len(elb.SecurityGroups)),
		})
	}

	table.Print(writer)
	if len(p.V1) > 0 || len(p.V2) > 0 {
		fmt.Println("* Indicates internal load balancer")
	}

	return nil
}

func (p *LoadBalancerPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.LoadBalancers)
}

func (p *LoadBalancerPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.LoadBalancers)
}
