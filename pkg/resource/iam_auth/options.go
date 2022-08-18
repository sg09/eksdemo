package iam_auth

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/eksctl"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
	"fmt"

	"github.com/spf13/cobra"
)

type IamAuthOptions struct {
	resource.CommonOptions
	eksctl.IamAuth
}

func NewOptions() (options *IamAuthOptions, flags cmd.Flags) {
	options = &IamAuthOptions{}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "arn",
				Description: "ARN of the IAM role or user",
				Validate: func(cmd *cobra.Command, args []string) error {
					t := template.TextTemplate{
						Template: options.Arn,
					}

					_, err := t.Render(options)
					if err != nil {
						return fmt.Errorf("failed to render ARN: %s", err)
					}

					return nil
				},
			},
			Option: &options.Arn,
		},
		&cmd.StringSliceFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "groups",
				Description: "Groups",
			},
			Option: &options.Groups,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "username",
				Description: "Kubernetes username to map IAM role or user",
			},
			Option: &options.Username,
		},
	}
	return
}
