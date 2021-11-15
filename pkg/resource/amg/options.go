package amg

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
	"strings"
)

type AmgOptions struct {
	resource.CommonOptions

	Auth          []string
	WorkspaceName string
}

func NewOptions() (options *AmgOptions, flags cmd.Flags) {
	options = &AmgOptions{
		CommonOptions: resource.CommonOptions{
			Name:                "amazon-managed-grafana",
			ClusterFlagDisabled: true,
			DeleteById:          true,
		},
	}

	flags = cmd.Flags{
		&cmd.StringSliceFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "auth",
				Description: "Authentication methods (aws_sso, saml)",
				Required:    true,
				Validate: func() error {
					for i, flag := range options.Auth {
						options.Auth[i] = strings.ToUpper(flag)
					}
					return nil
				},
			},
			Choices: []string{"aws_sso", "saml"},
			Option:  &options.Auth,
		},
	}

	return
}

func (o *AmgOptions) SetName(name string) {
	o.WorkspaceName = name
}
