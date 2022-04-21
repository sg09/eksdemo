package wordpress

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
)

type WordpressOptions struct {
	application.ApplicationOptions

	IngressHost       string
	StorageClass      string
	WordpressPassword string
}

func NewOptions() (options *WordpressOptions, flags cmd.Flags) {
	options = &WordpressOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:                 "wordpress",
			DisableServiceAccountFlag: true,
			DefaultVersion: &application.LatestPrevious{
				Latest:   "5.9.3",
				Previous: "5.9.3",
			},
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ingress-host",
				Description: "hostname for Ingress with TLS (requires ACM cert, AWS LB Controller and ExternalDNS)",
				Shorthand:   "I",
			},
			Option: &options.IngressHost,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "storage-class",
				Description: "StorageClass for WordPress and MariaDB Persistent Volumes",
			},
			Option: &options.StorageClass,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "wordpress-pass",
				Description: "WordPress admin password (required)",
				Required:    true,
			},
			Option: &options.WordpressPassword,
		},
	}
	return
}
