package application

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
	"fmt"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/spf13/cobra"
)

type Options interface {
	AddInstallFlags(*cobra.Command, cmd.Flags) cmd.Flags
	AddUninstallFlags(*cobra.Command, cmd.Flags, bool) cmd.Flags
	AssignCommonResourceOptions(*resource.Resource)
	Common() *ApplicationOptions
	KubeContext() string
	PreInstall() error
	PostInstall() error
}

type ApplicationOptions struct {
	Version string

	DefaultVersion
	DeleteDependencies        bool
	DisableNamespaceFlag      bool
	DisableServiceAccountFlag bool
	UsePrevious               bool

	Account        string
	ClusterName    string
	Namespace      string
	Region         string
	ServiceAccount string
	Cluster        *eks.Cluster
	kubeContext    string
}

type Action string

const Install Action = "install"
const Uninstall Action = "uninstall"

func (o *ApplicationOptions) AddInstallFlags(cobraCmd *cobra.Command, flags cmd.Flags) cmd.Flags {
	// Cluster flag has to be ordered before Version flag as it depends on the EKS cluster version
	flags = append(flags, o.NewClusterFlag(Install), o.NewVersionFlag(), o.NewUsePreviousFlag())

	if !o.DisableNamespaceFlag {
		flags = append(flags, o.NewNamespaceFlag(Install))
	}

	if !o.DisableServiceAccountFlag {
		flags = append(flags, o.NewServiceAccountFlag())
	}

	for _, f := range flags {
		f.AddFlagToCommand(cobraCmd)
	}

	return flags
}

func (o *ApplicationOptions) AddUninstallFlags(cobraCmd *cobra.Command, _ cmd.Flags, iamPolicy bool) cmd.Flags {
	commonFlags := cmd.Flags{
		o.NewClusterFlag(Uninstall),
		o.NewNamespaceFlag(Uninstall),
	}

	if iamPolicy {
		commonFlags = append(commonFlags, o.NewDeleteRoleFlag())
	}

	flags := commonFlags

	for _, f := range flags {
		f.AddFlagToCommand(cobraCmd)
	}

	return flags
}

func (o *ApplicationOptions) AssignCommonResourceOptions(res *resource.Resource) {
	r := res.Common()

	r.Account = aws.AccountId()
	r.Cluster = o.Cluster
	r.ClusterName = o.ClusterName
	r.KubeContext = o.kubeContext
	r.Namespace = o.Namespace
	r.Region = aws.Region()
	r.ServiceAccount = o.ServiceAccount
}

func (o *ApplicationOptions) Common() *ApplicationOptions {
	return o
}

func (o *ApplicationOptions) IrsaAnnotation() string {
	return fmt.Sprintf("eks.amazonaws.com/role-arn: arn:aws:iam::%s:role/eksdemo.%s.%s.%s",
		o.Account, o.ClusterName, o.Namespace, o.ServiceAccount)
}

func (o *ApplicationOptions) KubeContext() string {
	return o.kubeContext
}

func (o *ApplicationOptions) PreInstall() error {
	return nil
}

func (o *ApplicationOptions) PostInstall() error {
	return nil
}
