package service_linked_role

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResourceWithOptions(options *ServiceLinkedRoleOptions) *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "service-linked-role",
			Description: "Service Linked Role",
		},

		Manager: &Manager{},

		Options: options,
	}

	return res
}
