package resource

import (
	"eksdemo/pkg/cmd"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/spf13/cobra"
)

type Options interface {
	AddCreateFlags(*cobra.Command, cmd.Flags) cmd.Flags
	AddDeleteFlags(*cobra.Command, cmd.Flags) cmd.Flags
	AddGetFlags(*cobra.Command, cmd.Flags) cmd.Flags
	Common() *CommonOptions
	GetClusterName() string
	GetKubeContext() string
	PreCreate() error
	PreDelete() error
	SetName(string)
	Validate() error
}

type CommonOptions struct {
	Name                string
	ClusterFlagDisabled bool
	ClusterFlagOptional bool
	DeleteById          bool
	KubeContext         string
	NamespaceFlag       bool

	Account           string
	Cluster           *eks.Cluster
	ClusterName       string
	DryRun            bool
	KubernetesVersion string
	Namespace         string
	Region            string
	ServiceAccount    string
}

type Action string

const Create Action = "create"
const Delete Action = "delete"
const Get Action = "get"

func (o *CommonOptions) AddCreateFlags(cobraCmd *cobra.Command, flags cmd.Flags) cmd.Flags {
	flags = append(flags, o.NewDryRunFlag())

	if !o.ClusterFlagDisabled {
		flags = append(flags, o.NewClusterFlag(Create, true))
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

	if !o.ClusterFlagDisabled {
		flags = append(flags, o.NewClusterFlag(Delete, true))
	}

	if o.NamespaceFlag {
		flags = append(flags, o.NewNamespaceFlag(Delete))
	}

	for _, f := range flags {
		f.AddFlagToCommand(cobraCmd)
	}

	return flags
}

func (o *CommonOptions) AddGetFlags(cobraCmd *cobra.Command, _ cmd.Flags) cmd.Flags {
	flags := cmd.Flags{}

	if o.ClusterFlagOptional {
		flags = append(flags, o.NewClusterFlag(Get, false))
	} else if !o.ClusterFlagDisabled {
		flags = append(flags, o.NewClusterFlag(Get, true))
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

func (o *CommonOptions) GetKubeContext() string {
	return o.KubeContext
}

func (o *CommonOptions) PreCreate() error {
	return nil
}

func (o *CommonOptions) PreDelete() error {
	return nil
}

func (o *CommonOptions) SetName(name string) {
	o.Name = name
}

func (o *CommonOptions) Validate() error {
	return nil
}
