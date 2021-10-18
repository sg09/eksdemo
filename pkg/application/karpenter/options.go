package karpenter

import (
	"eksdemo/pkg/application"
	"fmt"
)

type KarpenterOptions struct {
	application.ApplicationOptions
}

func (o *KarpenterOptions) PostInstall() error {
	res := karpenterDefaultProvisioner()
	o.AssignCommonResourceOptions(res)
	fmt.Printf("Creating post-install resource: %s\n", res.Common().Name)

	return res.Create()
}
