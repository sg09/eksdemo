package cloudformation

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"fmt"

	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun               bool
	cloudformationClient *aws.CloudformationClient
	cloudformationGetter *Getter
}

func (m *Manager) Init() {
	if m.cloudformationClient == nil {
		m.cloudformationClient = aws.NewCloudformationClient()
	}
	m.cloudformationGetter = NewGetter(m.cloudformationClient)
}

func (m *Manager) Delete(options resource.Options) error {
	stackName := options.Common().Name

	_, err := m.cloudformationGetter.GetStacks(stackName)
	if err != nil {
		return err
	}

	fmt.Printf("Deleting Cloudformation stack %q\n", stackName)

	return m.cloudformationClient.DeleteStack(stackName)
}

func (m *Manager) Create(options resource.Options) error {
	return fmt.Errorf("feature not yet implemented")
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}
