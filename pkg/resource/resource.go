package resource

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/printer"
	"fmt"

	"github.com/spf13/cobra"
)

type Resource struct {
	cmd.Command
	cmd.Flags
	Options

	Getter
	Manager
}

func (r *Resource) Create() error {
	return r.Manager.Create(r.Options)
}

func (r *Resource) Delete() error {
	return r.Manager.Delete(r.Options)
}

func (r *Resource) NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     r.Command.Name + " NAME",
		Short:   r.Description,
		Long:    "Create " + r.Description,
		Aliases: r.Aliases,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := r.Flags.ValidateFlags(); err != nil {
				return err
			}

			if err := r.Options.Validate(); err != nil {
				return err
			}

			cmd.SilenceUsage = true
			r.SetName(args[0])

			if r.Common().DryRun {
				r.SetDryRun()
			}

			if err := r.PreCreate(); err != nil {
				return err
			}

			if r.Manager == nil {
				return fmt.Errorf("feature not yet implemented")
			}

			return r.Create()
		},
	}
	r.Flags = r.Options.AddCreateFlags(cmd, r.Flags)

	return cmd
}

func (r *Resource) NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     r.Command.Name + " NAME",
		Short:   r.Description,
		Long:    "Delete " + r.Description,
		Aliases: r.Aliases,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := r.Flags.ValidateFlags(); err != nil {
				return err
			}
			cmd.SilenceUsage = true

			r.SetName(args[0])

			if err := r.PreDelete(); err != nil {
				return err
			}

			if r.Manager == nil {
				return fmt.Errorf("feature not yet implemented")
			}

			return r.Delete()
		},
	}
	r.Flags = r.Options.AddDeleteFlags(cmd, r.Flags)

	return cmd
}

func (r *Resource) NewGetCmd() *cobra.Command {
	var output printer.Output

	cobraCmd := &cobra.Command{
		Use:     r.Command.Name + " NAME",
		Short:   r.Description,
		Long:    "Get " + r.Description,
		Aliases: r.Aliases,
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := r.Flags.ValidateFlags(); err != nil {
				return err
			}
			cmd.SilenceUsage = true

			if r.Getter == nil {
				return fmt.Errorf("feature not yet implemented")
			}

			name := ""
			if len(args) > 0 {
				name = args[0]
			}

			return r.Getter.Get(name, output, r.Common())
		},
	}
	cobraCmd.Flags().VarP(cmd.NewOutputFlag(&output), "output", "o", "output format (json|table|yaml)")

	r.Flags = r.Options.AddGetFlags(cobraCmd, r.Flags)

	return cobraCmd
}
