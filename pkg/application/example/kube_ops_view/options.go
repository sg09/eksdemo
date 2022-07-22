package kube_ops_view

import (
	"eksdemo/pkg/application"
	"eksdemo/pkg/cmd"
)

type KubeOpsViewOptions struct {
	application.ApplicationOptions

	Replicas int
}

func newOptions() (options *KubeOpsViewOptions, flags cmd.Flags) {
	options = &KubeOpsViewOptions{
		ApplicationOptions: application.ApplicationOptions{
			EnableIngress: true,
			Namespace:     "kube-ops-view",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "latest",
				Previous: "20.4.0",
			},
			DisableServiceAccountFlag: true,
		},
		Replicas: 1,
	}
	return
}
