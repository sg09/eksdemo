package network_interface

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"fmt"
	"io"
	"strconv"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type NetworkInterfacePrinter struct {
	networkInterfaces []*ec2.NetworkInterface
}

func NewPrinter(networkInterfaces []*ec2.NetworkInterface) *NetworkInterfacePrinter {
	return &NetworkInterfacePrinter{networkInterfaces}
}

func (p *NetworkInterfacePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Id", "Instance Id or...", "Private IPv4", "IPs", "SGs", "Subnet"})

	for _, eni := range p.networkInterfaces {
		id := aws.StringValue(eni.NetworkInterfaceId)
		instanceId := ""

		if eni.Attachment == nil {
			instanceId = "detached"
		} else {
			if aws.Int64Value(eni.Attachment.DeviceIndex) == 0 {
				id = "*" + id
			}

			if aws.StringValue(eni.Attachment.InstanceId) != "" {
				instanceId = aws.StringValue(eni.Attachment.InstanceId)
			} else if aws.StringValue(eni.InterfaceType) != "interface" {
				instanceId = aws.StringValue(eni.InterfaceType)
			} else if aws.StringValue(eni.Attachment.InstanceOwnerId) == "amazon-elb" {
				instanceId = "elb"
			}

			// Identify EKS Control Plane cross-account ENI
			// This identifies any interface ENI that is attached to EC2 but owned by a different account
			_, err := strconv.Atoi(aws.StringValue(eni.Attachment.InstanceOwnerId))
			if aws.StringValue(eni.InterfaceType) == "interface" && aws.StringValue(eni.Attachment.Status) == "attached" {
				if err == nil && aws.StringValue(eni.Attachment.InstanceOwnerId) != aws.AccountId() {
					instanceId = "eks_control_plane"
				}
			}
		}

		table.AppendRow([]string{
			id,
			instanceId,
			aws.StringValue(eni.PrivateIpAddress),
			strconv.Itoa(len(eni.PrivateIpAddresses)),
			strconv.Itoa(len(eni.Groups)),
			aws.StringValue(eni.SubnetId),
		})
	}

	table.Print(writer)
	if len(p.networkInterfaces) > 0 {
		fmt.Println("* Indicates Primary network interface")
	}

	return nil
}

func (p *NetworkInterfacePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.networkInterfaces)
}

func (p *NetworkInterfacePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.networkInterfaces)
}
