package cluster

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/application/aws_lb"
	"eksdemo/pkg/application/cluster_autoscaler"
	"eksdemo/pkg/application/external_dns"
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/resource/nodegroup"
	"eksdemo/pkg/template"
)

type ClusterOptions struct {
	resource.CommonOptions
	*nodegroup.NodegroupOptions

	Fargate bool
	NoRoles bool

	appsForIrsa  []*application.Application
	IrsaTemplate *template.TextTemplate
	IrsaRoles    []*resource.Resource
}

func addOptions(res *resource.Resource) *resource.Resource {
	ngOptions, ngFlags := nodegroup.NewOptions()
	ngOptions.DesiredCapacity = 2
	ngOptions.MinSize = 2
	ngOptions.NodegroupName = "main"

	options := &ClusterOptions{
		CommonOptions: resource.CommonOptions{
			DisableClusterFlag: true,
			KubernetesVersion:  "1.21",
		},

		NodegroupOptions: ngOptions,
		NoRoles:          false,

		appsForIrsa: []*application.Application{
			aws_lb.NewApp(),
			cluster_autoscaler.NewApp(),
			external_dns.NewApp(),
		},
		IrsaTemplate: &template.TextTemplate{
			Template: irsa.EksctlTemplate,
		},
	}

	res.Options = options

	flags := cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "version",
				Description: "Kubernetes version",
				Shorthand:   "v",
			},
			Choices: []string{"1.21", "1.20", "1.19", "1.18", "1.17"},
			Option:  &options.KubernetesVersion,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "no-roles",
				Description: "don't create IAM roles",
			},
			Option: &options.NoRoles,
		},
	}

	res.Flags = append(ngFlags, flags...)

	return res
}

func (o *ClusterOptions) PreCreate() error {
	o.Account = aws.AccountId()
	o.NodegroupOptions.KubernetesVersion = o.KubernetesVersion

	// Add ClusterName to the Applications used to pre-create IRSA roles and
	// Rreate IRSA Resources that can be used to render the role for eksctl
	for _, app := range o.appsForIrsa {
		app.Common().ClusterName = o.ClusterName

		for _, res := range app.Dependencies {
			if res.Name != "irsa" {
				continue
			}
			app.AssignCommonResourceOptions(res)
			o.IrsaRoles = append(o.IrsaRoles, res)
		}
	}

	return o.NodegroupOptions.PreCreate()
}

func (o *ClusterOptions) SetName(name string) {
	o.ClusterName = name
	o.Region = aws.Region()
}
