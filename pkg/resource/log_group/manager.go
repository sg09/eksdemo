package log_group

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"fmt"

	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun               bool
	cloudwatchlogsClient *aws.CloudWatchLogsClient
}

func (m *Manager) Init() {
	if m.cloudwatchlogsClient == nil {
		m.cloudwatchlogsClient = aws.NewCloudWatchLogsClient()
	}
}

func (m *Manager) Create(options resource.Options) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) Delete(options resource.Options) error {
	if err := m.cloudwatchlogsClient.DeleteLogGroup(options.Common().Name); err != nil {
		return aws.FormatError(err)
	}
	fmt.Println("Log group deleted successfully")

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}
