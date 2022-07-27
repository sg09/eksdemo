package cluster

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/application/autoscaling/cluster_autoscaler"
	"eksdemo/pkg/application/autoscaling/karpenter"
	"eksdemo/pkg/application/aws_lb"
	"eksdemo/pkg/application/external_dns"
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/cloudformation"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/resource/nodegroup"
	"eksdemo/pkg/template"
	"fmt"
	"strings"
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
	ngOptions, ngFlags, _ := nodegroup.NewOptions()
	ngOptions.DesiredCapacity = 2
	ngOptions.NodegroupName = "main"

	options := &ClusterOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
			KubernetesVersion:   "1.22",
		},

		NodegroupOptions: ngOptions,
		NoRoles:          false,

		appsForIrsa: []*application.Application{
			aws_lb.NewApp(),
			cluster_autoscaler.NewApp(),
			external_dns.NewApp(),
			karpenter.NewApp(),
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
			Choices: []string{"1.22", "1.21", "1.20", "1.19"},
			Option:  &options.KubernetesVersion,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "fargate",
				Description: "create a Fargate profile",
			},
			Option: &options.Fargate,
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

	// For apps we want to pre-create IRSA for, find the IRSA dependency
	for _, app := range o.appsForIrsa {
		for _, res := range app.Dependencies {
			if res.Name != "irsa" {
				continue
			}
			// Populate the IRSA Resource with data (Cluster, Namespace, ServiceAccount)
			app.Common().ClusterName = o.ClusterName
			app.AssignCommonResourceOptions(res)
			res.SetName(app.Common().ServiceAccount)

			o.IrsaRoles = append(o.IrsaRoles, res)
		}
	}

	return o.NodegroupOptions.PreCreate()
}

func (o *ClusterOptions) PreDelete() error {
	getter := cloudformation.Getter{}
	stacks, err := getter.GetStacksByCluster(o.ClusterName, "")
	if err != nil {
		return err
	}

	for _, stack := range stacks {
		stackName := aws.StringValue(stack.StackName)
		if strings.HasPrefix(stackName, "eksdemo-") {
			fmt.Printf("Deleting Cloudformation stack %q\n", stackName)
			err := aws.CloudFormationDeleteStack(stackName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (o *ClusterOptions) SetName(name string) {
	o.ClusterName = name
	o.Region = aws.Region()
}
