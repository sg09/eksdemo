package irsa

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
	"fmt"
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

	res.Flags = cmd.Flags{
		&cmd.StringSliceFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "attach-arns",
				Description: "ARNs",
				Validate: func() error {
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
				Validate: func() error {
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

func (o *IrsaOptions) SetName(name string) {
	o.ServiceAccount = name
}

func (o *IrsaOptions) Validate() error {
	if o.PolicyType == None {
		return fmt.Errorf("a single policy type must be used")
	}
	return nil
}
