package amg

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
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
			ArgumentOptional:    true,
			ClusterFlagDisabled: true,
			DeleteByIdFlag:      true,
		},
	}

	flags = cmd.Flags{
		&cmd.StringSliceFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "auth",
				Description: "Authentication methods (aws_sso, saml)",
				Required:    true,
				Validate: func(cmd *cobra.Command, args []string) error {
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

func (o *AmgOptions) iamRoleName() string {
	return fmt.Sprintf("eksdemo.amg.%s", o.WorkspaceName)
}
