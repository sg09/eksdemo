package irsa

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/spf13/cobra"
)

type IrsaOptions struct {
	resource.CommonOptions

	PolicyType
	Policy []string

	// Used for flags
	WellKnownPolicy   string
	PolicyARNs        []string
	PolicyDocTemplate template.Template
}

type PolicyType int

const (
	None PolicyType = iota
	WellKnown
	PolicyARNs
	PolicyDocument
)

func addOptions(res *resource.Resource) *resource.Resource {
	options := &IrsaOptions{
		CommonOptions: resource.CommonOptions{
			Namespace:     "default",
			NamespaceFlag: true,
		},
	}

	res.Options = options

	res.CreateFlags = cmd.Flags{
		&cmd.StringSliceFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "attach-arns",
				Description: "ARNs",
				Validate: func(cmd *cobra.Command, args []string) error {
					if len(options.PolicyARNs) == 0 {
						return nil
					}

					if len(options.Policy) > 0 {
						return fmt.Errorf("can only use one policy flag")
					}

					options.PolicyType = PolicyARNs
					options.Policy = options.PolicyARNs

					return nil
				},
			},
			Option: &options.PolicyARNs,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "well-known",
				Description: "eksctl well known policy",
				Validate: func(cmd *cobra.Command, args []string) error {
					if options.WellKnownPolicy == "" {
						return nil
					}

					if len(options.Policy) > 0 {
						return fmt.Errorf("can only use one policy flag")
					}

					options.PolicyType = WellKnown
					options.Policy = []string{options.WellKnownPolicy}

					return nil
				},
			},
			Option: &options.WellKnownPolicy,
		},
	}

	return res
}

func (o *IrsaOptions) ClusterOIDCProvider() (string, error) {
	issuer := aws.ToString(o.Cluster.Identity.Oidc.Issuer)

	slices := strings.Split(issuer, "//")
	if len(slices) < 2 {
		return "", fmt.Errorf("failed to parse Cluster OIDC Provider URL")
	}

	return slices[1], nil
}

func (o *IrsaOptions) IrsaAnnotation() string {
	return fmt.Sprintf("eks.amazonaws.com/role-arn: arn:%s:iam::%s:role/%s", o.Partition, o.Account, o.RoleName())
}

func (o *IrsaOptions) IsPolicyDocument(t PolicyType) bool {
	return t == PolicyDocument
}

func (o *IrsaOptions) IsPolicyARN(t PolicyType) bool {
	return t == PolicyARNs
}

func (o *IrsaOptions) IsWellKnownPolicy(t PolicyType) bool {
	return t == WellKnown
}

func (o *IrsaOptions) PreDelete() error {
	o.PolicyType = WellKnown
	o.Policy = []string{"autoScaler"}
	return nil
}

func (o *IrsaOptions) RoleName() string {
	return o.TruncateUnique(64, fmt.Sprintf("eksdemo.%s.%s.%s", o.ClusterName, o.Namespace, o.ServiceAccount))
}

func (o *IrsaOptions) SetName(name string) {
	o.ServiceAccount = name
}
