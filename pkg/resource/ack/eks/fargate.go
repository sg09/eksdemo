package eks

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/kubernetes"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
	"eksdemo/pkg/util"
)

type FargateProfileOptions struct {
	resource.CommonOptions
	FargateNamespace string
	Subnets          []string
}

func NewFargateProfileResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "eks-fargate-profile",
			Description: "Fargate Profile",
			Aliases:     []string{"eks-fargate", "fargate-profile", "fargate", "fp"},
			Args:        []string{"NAME"},
		},

		Manager: &kubernetes.ResourceManager{
			Template: &template.TextTemplate{
				Template: subnetYamlTemplate,
			},
		},
	}

	options := &FargateProfileOptions{
		CommonOptions: resource.CommonOptions{
			Name:          "ack-eks-fargate-profile",
			Namespace:     "default",
			NamespaceFlag: true,
		},
		FargateNamespace: "default",
	}

	flags := cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "fargate-namespace",
				Description: "namespace selector to run pods on fargate",
			},
			Option: &options.FargateNamespace,
		},
		&cmd.StringSliceFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "subnets",
				Description: "subnets for fargate pods (defaults to all private subnets)",
			},
			Option: &options.Subnets,
		},
	}

	res.Options = options
	res.CreateFlags = flags

	return res
}

func (o *FargateProfileOptions) PreCreate() error {
	if len(o.Subnets) > 0 {
		return nil
	}

	subnets, err := util.GetPrivateSubnets(o.ClusterName)
	if err != nil {
		return err
	}

	o.Subnets = subnets
	return nil
}

const subnetYamlTemplate = `
---
apiVersion: eks.services.k8s.aws/v1alpha1
kind: FargateProfile
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  clusterName: {{ .ClusterName }}
  name: {{ .Name }}
  podExecutionRoleARN: arn:aws:iam::{{ .Account }}:role/eksdemo.{{ .ClusterName }}.fargate-pod-execution-role
  selectors:
  - namespace: {{ .FargateNamespace }}
  subnets:
{{- range .Subnets }}
  - {{ . }}
{{- end }}
...
`
