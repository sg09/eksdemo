package ec2

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/kubernetes"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
)

type SubnetOptions struct {
	resource.CommonOptions
	CidrBlock string
	VpcId     string
}

func NewSubnetResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "subnet",
			Description: "EC2 Subnet",
			Aliases:     []string{"subnets"},
			Args:        []string{"NAME"},
		},

		Manager: &kubernetes.ResourceManager{
			Template: &template.TextTemplate{
				Template: subnetYamlTemplate,
			},
		},
	}

	options := &SubnetOptions{
		CommonOptions: resource.CommonOptions{
			Name:          "ack-ec2-subnet",
			Namespace:     "default",
			NamespaceFlag: true,
		},
	}

	flags := cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "cidr",
				Description: "ipv4 network range for the VPC, in CIDR notation (required)",
				Required:    true,
			},
			Option: &options.CidrBlock,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "vpc-id",
				Description: "ip of the VPC to create the subnet (defaults to cluster VPC)",
			},
			Option: &options.VpcId,
		},
	}

	res.Options = options
	res.Flags = flags

	return res
}

func (o *SubnetOptions) PreCreate() error {
	if o.VpcId == "" {
		o.VpcId = aws.StringValue(o.Cluster.ResourcesVpcConfig.VpcId)
	}
	return nil
}

const subnetYamlTemplate = `
---
apiVersion: ec2.services.k8s.aws/v1alpha1
kind: Subnet
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  cidrBlock: {{ .CidrBlock }}
  tagSpecifications:
  - resourceType: subnet
    tags:
    - key: Name
      value: {{ .Name }}
  vpcID: {{ .VpcId }}
...
`
