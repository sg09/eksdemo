package iam_auth

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/eksctl"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	options, flags := NewOptions()
	res := NewResourceWithOptions(options)
	res.CreateFlags = flags

	return res
}

func NewResourceWithOptions(options *IamAuthOptions) *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "iam-auth",
			Description: "IAM User or Role for Authentication",
		},

		Manager: &eksctl.ResourceManager{
			Resource: "iamidentitymapping",
			IamAuth:  &options.IamAuth,
		},
	}

	res.Options = options

	return res
}
