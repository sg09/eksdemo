package ec2

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/kubernetes"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
)

type VpcOptions struct {
	resource.CommonOptions
	CidrBlock string
}

func NewVpcResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "vpc",
			Description: "Virtual Private Cloud (VPC)",
			Aliases:     []string{"vpcs"},
			Args:        []string{"NAME"},
		},

		Manager: &kubernetes.ResourceManager{
			Template: &template.TextTemplate{
				Template: vpcYamlTemplate,
			},
		},
	}

	options := &VpcOptions{
		CommonOptions: resource.CommonOptions{
			Name:          "ack-ec2-vpc",
			Namespace:     "default",
			NamespaceFlag: true,
		},
	}

	flags := cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "cidr",
				Description: "ipv4 network range for the VPC, in CIDR notation",
				Required:    true,
			},
			Option: &options.CidrBlock,
		},
	}

	res.Options = options
	res.CreateFlags = flags

	return res
}

const vpcYamlTemplate = `
---
apiVersion: ec2.services.k8s.aws/v1alpha1
kind: VPC
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  cidrBlock: {{ .CidrBlock }}
  tagSpecifications:
  - resourceType: vpc
    tags:
    - key: Name
      value: {{ .Name }}
...
`
