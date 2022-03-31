package eks_workshop

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
)

type EksWorkshopOptions struct {
	application.ApplicationOptions

	CrystalReplicas  int
	FrontendReplicas int
	IngressHost      string
	NodeReplicas     int
}

func NewOptions() (options *EksWorkshopOptions, flags cmd.Flags) {
	options = &EksWorkshopOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:                 "eks-workshop",
			DisableServiceAccountFlag: true,
			DisableVersionFlag:        true,
		},
		CrystalReplicas:  3,
		FrontendReplicas: 3,
		NodeReplicas:     3,
	}

	flags = cmd.Flags{
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "crystal-replicas",
				Description: "number of replicas for the Crystal deployment",
			},
			Option: &options.CrystalReplicas,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "frontend-replicas",
				Description: "number of replicas for the Frontend deployment",
			},
			Option: &options.FrontendReplicas,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ingress-host",
				Description: "hostname for Ingress with TLS (requires ACM cert, AWS LB Controller and ExternalDNS)",
				Shorthand:   "I",
			},
			Option: &options.IngressHost,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "nodejs-replicas",
				Description: "number of replicas for the Node.js deployment",
			},
			Option: &options.NodeReplicas,
		},
	}
	return
}
