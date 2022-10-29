package load_balancer

import (
	"eksdemo/pkg/printer"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
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

	for _, lb := range p.V1 {
		age := durafmt.ParseShort(time.Since(aws.ToTime(lb.CreatedTime)))

		name := aws.ToString(lb.LoadBalancerName)
		if aws.ToString(lb.Scheme) == "internal" {
			name = "*" + name
		}

		table.AppendRow([]string{
			age.String(),
			"-",
			name,
			"CLB",
			"ipv4",
			strconv.Itoa(len(lb.AvailabilityZones)),
			strconv.Itoa(len(lb.SecurityGroups)),
		})
	}

	for _, elb := range p.V2 {
		age := durafmt.ParseShort(time.Since(aws.ToTime(elb.CreatedTime)))

		elbType := "unknown"
		switch string(elb.Type) {
		case "application":
			elbType = "ALB"
		case "network":
			elbType = "NLB"
		case "gateway":
			elbType = "GWLB"
		}

		table.AppendRow([]string{
			age.String(),
			string(elb.State.Code),
			aws.ToString(elb.LoadBalancerName),
			elbType,
			string(elb.IpAddressType),
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
