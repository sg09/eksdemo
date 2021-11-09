package cloudformation

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"fmt"
)

type Manager struct {
	DryRun bool
}

func (m *Manager) Delete(options resource.Options) error {
	stackName := options.Common().Name

	fmt.Printf("Deleting Cloudformation stack %q\n", stackName)

	return aws.CloudFormationDeleteStack(stackName)
}

func (m *Manager) Create(options resource.Options) error {
	return fmt.Errorf("feature not yet implemented")
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}
