package kubecost

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
	"fmt"
)

type KubecostOptions struct {
	application.ApplicationOptions
	AdminPassword string
}

func newOptions() (options *KubecostOptions, flags cmd.Flags) {
	options = &KubecostOptions{
		ApplicationOptions: application.ApplicationOptions{
			EnableIngress:  true,
			Namespace:      "kubecost",
			ServiceAccount: "kubecost-cost-analyzer",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.95.0",
				Latest:        "1.95.0",
				PreviousChart: "1.94.3",
				Previous:      "1.94.3",
			},
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "admin-pass",
				Description: "password for admin user (only valid with --ingress-class=nginx)",
				Shorthand:   "P",
				Validate: func() error {
					if options.AdminPassword != "" && options.IngressClass != "nginx" {
						return fmt.Errorf("%q flag can only be used with %q)", "admin-pass", "--ingress-class=nginx")
					}
					return nil
				},
			},
			Option: &options.AdminPassword,
		},
	}
	return
}

func (o *KubecostOptions) PostInstall(name string, _ []*resource.Resource) error {
	if o.IngressClass == "nginx" && o.AdminPassword != "" {
		return o.ApplicationOptions.PostInstall(name, []*resource.Resource{nginxSecret(o.AdminPassword)})
	}
	return nil
}
