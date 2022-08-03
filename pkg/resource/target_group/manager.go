package target_group

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"fmt"

	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun bool
	Getter
}

func (m *Manager) Create(options resource.Options) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) Delete(options resource.Options) error {
	name := options.Common().Name

	tg, err := m.Getter.GetTargetGroupByName(name)
	if err != nil {
		return err
	}

	if err := aws.ELBDeleteTargetGroup(aws.StringValue(tg.TargetGroupArn)); err != nil {
		return err
	}
	fmt.Printf("Target Group %q deleted\n", name)

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}
