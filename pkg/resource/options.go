package resource

import (
	"eksdemo/pkg/cmd"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/spf13/cobra"
)

type Options interface {
	AddCreateFlags(*cobra.Command, cmd.Flags) cmd.Flags
	AddDeleteFlags(*cobra.Command, cmd.Flags) cmd.Flags
	Common() *CommonOptions
	GetClusterName() string
	SetName(string)
	PreCreate() error
	PrepForDelete()
	Validate() error
}

type CommonOptions struct {
	Account           string
	ClusterName       string
	Name              string
	Namespace         string
	Region            string
	KubernetesVersion string

	DisableClusterFlag bool
	NamespaceFlag      bool

	Cluster *eks.Cluster
}

type Action string

const Create Action = "create"
const Delete Action = "delete"

func (o *CommonOptions) AddCreateFlags(cobraCmd *cobra.Command, flags cmd.Flags) cmd.Flags {
	if !o.DisableClusterFlag {
		flags = append(flags, o.NewClusterFlag(Create))
	}

	if o.NamespaceFlag {
		flags = append(flags, o.NewNamespaceFlag(Create))
	}

	for _, f := range flags {
		f.AddFlagToCommand(cobraCmd)
	}

	return flags
}

func (o *CommonOptions) AddDeleteFlags(cobraCmd *cobra.Command, _ cmd.Flags) cmd.Flags {
	flags := cmd.Flags{}

	if !o.DisableClusterFlag {
		flags = append(flags, o.NewClusterFlag(Delete))
	}

	if o.NamespaceFlag {
		flags = append(flags, o.NewNamespaceFlag(Delete))
	}

	for _, f := range flags {
		f.AddFlagToCommand(cobraCmd)
	}

	return flags
}

func (o *CommonOptions) Common() *CommonOptions {
	return o
}

func (o *CommonOptions) GetClusterName() string {
	return o.ClusterName
}

func (o *CommonOptions) PreCreate() error {
	return nil
}

func (o *CommonOptions) PrepForDelete() {}

func (o *CommonOptions) SetName(name string) {
	o.Name = name
}

func (o *CommonOptions) Validate() error {
	return nil
}
