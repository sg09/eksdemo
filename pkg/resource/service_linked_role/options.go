package service_linked_role

import (
	"eksdemo/pkg/resource"
)

type ServiceLinkedRoleOptions struct {
	resource.CommonOptions

	RoleName    string
	ServiceName string
}
