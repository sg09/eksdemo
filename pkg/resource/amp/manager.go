package amp

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"fmt"
)

type Manager struct{}

func (m *Manager) Create(options resource.Options) error {
	amp, ok := options.(*AmpOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to AmpOptions")
	}

	fmt.Printf("Creating AMP with Alias: %s...", amp.Alias)
	result, err := aws.AmpCreateWorkspace(amp.Alias)
	if err != nil {
		return err
	}
	fmt.Printf("done\nCreated AMP Workspace Id: %s\n", *result.WorkspaceId)

	return nil
}

func (m *Manager) Delete(options resource.Options) error {
	amp, ok := options.(*AmpOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to AmpOptions")
	}

	err := aws.AmpDeleteWorkspace(amp.Alias)
	if err != nil {
		return err
	}
	fmt.Printf("AMP %q deleted\n", amp.Alias)

	return nil
}

func (m *Manager) SetDryRun() {}
