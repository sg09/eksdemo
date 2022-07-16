package s3_bucket

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

type BucketOptions struct {
	resource.CommonOptions

	BucketName string
}

func newOptions() (options *BucketOptions, flags cmd.Flags) {
	options = &BucketOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}
	return
}
