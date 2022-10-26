package volume

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"fmt"

	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun    bool
	ec2Client *aws.EC2Client
}

func (m *Manager) Init() {
	if m.ec2Client == nil {
		m.ec2Client = aws.NewEC2Client()
	}
}

func (m *Manager) Create(options resource.Options) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) Delete(options resource.Options) error {
	if err := m.ec2Client.DeleteVolume(options.Common().Name); err != nil {
		return err
	}
	fmt.Println("Volume deleted successfully")

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}
