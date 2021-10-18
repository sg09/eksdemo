package application

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"fmt"
)

func (a *Application) CreateDependencies() error {
	if len(a.Dependencies) > 0 {
		fmt.Printf("Creating %d dependencies for %s\n", len(a.Dependencies), a.Name)
	}

	for _, res := range a.Dependencies {
		fmt.Printf("Creating dependency: %s\n", res.Common().Name)

		a.AssignCommonResourceOptions(res)

		if err := res.Create(); err != nil {
			return err
		}
	}
	return nil
}

func (a *Application) DeleteDependencies() error {
	if len(a.Dependencies) > 0 {
		fmt.Printf("Deleting %d dependencies for %s\n", len(a.Dependencies), a.Name)
	}

	for _, res := range a.Dependencies {
		fmt.Printf("Deleting dependency: %s\n", res.Common().Name)

		a.AssignCommonResourceOptions(res)

		if err := res.Delete(); err != nil {
			return err
		}
	}
	return nil
}

func (a *Application) AssignCommonResourceOptions(res *resource.Resource) {
	r := res.Common()

	r.Account = aws.AccountId()
	r.Cluster = a.Common().Cluster
	r.ClusterName = a.Common().ClusterName
	r.KubeContext = a.KubeContext()
	r.Namespace = a.Common().Namespace
	r.Region = aws.Region()
	r.ServiceAccount = a.Common().ServiceAccount
}
