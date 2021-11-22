package amp

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"fmt"
)

type Manager struct{}

func (m *Manager) Create(options resource.Options) error {
	ampOptions, ok := options.(*AmpOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to AmpOptions")
	}

	ampGetter := Getter{}
	workspace, err := ampGetter.GetAmpByAlias(ampOptions.Alias)
	if err != nil {
		if _, ok := err.(resource.NotFoundError); !ok {
			// Return an error if it's anything other than resource not found
			return err
		}
	}

	if workspace != nil {
		fmt.Printf("AMP Workspace Alias %q already exists\n", ampOptions.Alias)
		return nil
	}

	fmt.Printf("Creating AMP Workspace Alias: %s...", ampOptions.Alias)
	result, err := aws.AmpCreateWorkspace(ampOptions.Alias)
	if err != nil {
		return err
	}
	fmt.Printf("done\nCreated AMP Workspace Id: %s\n", *result.WorkspaceId)

	return nil
}

func (m *Manager) Delete(options resource.Options) error {
	ampOptions, ok := options.(*AmpOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to AmpOptions")
	}

	id := options.Common().Id

	if id == "" {
		ampGetter := Getter{}
		amp, err := ampGetter.GetAmpByAlias(ampOptions.Alias)
		if err != nil {
			if _, ok := err.(resource.NotFoundError); ok {
				fmt.Printf("AMP Workspace Alias %q does not exist\n", ampOptions.Alias)
				return nil
			}
			return err
		}
		id = aws.StringValue(amp.WorkspaceId)
	}

	err := aws.AmpDeleteWorkspace(id)
	if err != nil {
		return err
	}
	fmt.Printf("AMP Workspace Id %q deleting...\n", id)

	return nil
}

func (m *Manager) SetDryRun() {}
