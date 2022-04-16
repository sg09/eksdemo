package organization

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"fmt"

	"github.com/spf13/cobra"
)

type Manager struct{}

func (m *Manager) Create(options resource.Options) error {
	result, err := aws.OrgsCreateOrganization()
	if err != nil {
		return err
	}
	fmt.Printf("Created AWS Organization: %s\n", *result.Id)

	return nil
}

func (m *Manager) Delete(options resource.Options) error {
	err := aws.OrgsDeleteOrganization()
	if err != nil {
		return err
	}
	fmt.Println("AWS Organization deleted")

	return nil
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("update not supported")
}

func (m *Manager) SetDryRun() {}
