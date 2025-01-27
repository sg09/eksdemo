package cluster

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/application/autoscaling/cluster_autoscaler"
	"eksdemo/pkg/application/autoscaling/karpenter"
	"eksdemo/pkg/application/aws_lb_controller"
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

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
)

type ClusterOptions struct {
	resource.CommonOptions
	*nodegroup.NodegroupOptions

	Fargate bool
	IPv6    bool
	NoRoles bool
	Private bool

	appsForIrsa  []*application.Application
	IrsaTemplate *template.TextTemplate
	IrsaRoles    []*resource.Resource
}

func addOptions(res *resource.Resource) *resource.Resource {
	ngOptions, ngFlags, _ := nodegroup.NewOptions()

	options := &ClusterOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
			KubernetesVersion:   "1.24",
		},

		NodegroupOptions: ngOptions,
		NoRoles:          false,

		appsForIrsa: []*application.Application{
			aws_lb_controller.NewApp(),
			cluster_autoscaler.NewApp(),
			external_dns.NewApp(),
			karpenter.NewApp(),
		},
		IrsaTemplate: &template.TextTemplate{
			Template: irsa.EksctlTemplate,
		},
	}

	ngOptions.CommonOptions = options.Common()
	ngOptions.DesiredCapacity = 2
	ngOptions.NodegroupName = "main"

	res.Options = options

	flags := cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "version",
				Description: "Kubernetes version",
				Shorthand:   "v",
			},
			Choices: []string{"1.24", "1.23", "1.22", "1.21"},
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
				Name:        "ipv6",
				Description: "use IPv6 networking",
			},
			Option: &options.IPv6,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "no-roles",
				Description: "don't create IAM roles",
			},
			Option: &options.NoRoles,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "private",
				Description: "private cluster (includes ECR, S3, and other VPC endpoints)",
			},
			Option: &options.Private,
		},
	}

	res.CreateFlags = append(ngFlags, flags...)

	return res
}

func (o *ClusterOptions) PreCreate() error {
	o.Account = aws.AccountId()
	o.Partition = aws.Partition()
	o.Region = aws.Region()
	o.NodegroupOptions.KubernetesVersion = o.KubernetesVersion

	// For apps we want to pre-create IRSA for, find the IRSA dependency
	for _, app := range o.appsForIrsa {
		for _, res := range app.Dependencies {
			if res.Name != "irsa" {
				continue
			}
			// Populate the IRSA Resource with Account, Cluster, Namespace, Partition, Region, ServiceAccount
			app.Common().Account = o.Account
			app.Common().ClusterName = o.ClusterName
			app.Common().Region = o.Region
			app.Common().Partition = o.Partition
			app.AssignCommonResourceOptions(res)
			res.SetName(app.Common().ServiceAccount)

			o.IrsaRoles = append(o.IrsaRoles, res)
		}
	}

	return o.NodegroupOptions.PreCreate()
}

func (o *ClusterOptions) PreDelete() error {
	o.Region = aws.Region()

	cloudformationClient := aws.NewCloudformationClient()
	stacks, err := cloudformation.NewGetter(cloudformationClient).GetStacksByCluster(o.ClusterName, "")
	if err != nil {
		return err
	}

	for _, stack := range stacks {
		stackName := awssdk.ToString(stack.StackName)
		if strings.HasPrefix(stackName, "eksdemo-") {
			fmt.Printf("Deleting Cloudformation stack %q\n", stackName)
			err := cloudformationClient.DeleteStack(stackName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (o *ClusterOptions) SetName(name string) {
	o.ClusterName = name
}
