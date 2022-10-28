package fargate_profile

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

type FargateProfileOptions struct {
	resource.CommonOptions

	Namespaces []string
}

func NewOptions() (options *FargateProfileOptions, flags cmd.Flags) {
	options = &FargateProfileOptions{
		Namespaces: []string{"default"},
	}

	flags = cmd.Flags{
		&cmd.StringSliceFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "namespaces",
				Description: "namespaces to select pods from",
			},
			Option: &options.Namespaces,
		},
	}
	return
}
