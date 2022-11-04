package kms_key

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

type KmsKeyOptions struct {
	resource.CommonOptions
}

func newOptions() (options *KmsKeyOptions, getFlags cmd.Flags) {
	options = &KmsKeyOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	getFlags = cmd.Flags{}

	return
}
