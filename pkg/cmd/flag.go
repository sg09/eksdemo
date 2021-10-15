package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type Flag interface {
	AddFlagToCommand(*cobra.Command)
	ValidateFlag() error
	GetName() string
}

type Flags []Flag

type CommandFlag struct {
	Name        string
	Description string
	Shorthand   string
	Required    bool
	Validate    func() error
}

type BoolFlag struct {
	CommandFlag
	Option *bool
}

type IntFlag struct {
	CommandFlag
	Option *int
}

type StringFlag struct {
	CommandFlag
	Choices []string
	Option  *string
}

type StringSliceFlag struct {
	CommandFlag
	Option *[]string
}

// Command flag methods
func (f *CommandFlag) GetName() string {
	return f.Name
}

// Boolean flag methods
func (f *BoolFlag) AddFlagToCommand(cmd *cobra.Command) {
	cmd.Flags().BoolVar(f.Option, f.Name, *f.Option, f.Description)
	if f.Required {
		cmd.MarkFlagRequired(f.Name)
	}
}

func (f *BoolFlag) ValidateFlag() error {
	if f.Validate == nil {
		return nil
	}
	return f.Validate()
}

// Int flag methods
func (f *IntFlag) AddFlagToCommand(cmd *cobra.Command) {
	cmd.Flags().IntVarP(f.Option, f.Name, f.Shorthand, *f.Option, f.Description)
	if f.Required {
		cmd.MarkFlagRequired(f.Name)
	}
}

func (f *IntFlag) ValidateFlag() error {
	if f.Validate == nil {
		return nil
	}
	return f.Validate()
}

// String flag methods
func (f *StringFlag) AddFlagToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVarP(f.Option, f.Name, f.Shorthand, *f.Option, f.Description)
	if f.Required {
		cmd.MarkFlagRequired(f.Name)
	}
}

func (f *StringFlag) ValidateFlag() error {
	if len(f.Choices) > 0 {
		found := false

		for _, choice := range f.Choices {
			if strings.EqualFold(choice, *f.Option) {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("--%s must be one of: %s", f.Name, strings.Join(f.Choices, ", "))
		}
	}

	if f.Validate != nil {
		return f.Validate()
	}

	return nil
}

// StringSlice flag methods
func (f *StringSliceFlag) AddFlagToCommand(cmd *cobra.Command) {
	cmd.Flags().StringSliceVarP(f.Option, f.Name, f.Shorthand, *f.Option, f.Description)
	if f.Required {
		cmd.MarkFlagRequired(f.Name)
	}
}

func (f *StringSliceFlag) ValidateFlag() error {
	if f.Validate == nil {
		return nil
	}
	return f.Validate()
}

// Flags (list of flags) methods
func (f Flags) ValidateFlags() error {
	for _, flag := range f {
		if err := flag.ValidateFlag(); err != nil {
			return err
		}
	}
	return nil
}

func (f Flags) Remove(name string) Flags {
	for i, flag := range f {
		if flag.GetName() == name {
			f[i] = f[len(f)-1]
			f = f[:len(f)-1]
			break
		}
	}
	return f
}
