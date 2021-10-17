package application

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
)

func (a *Application) CreateDependencies() error {
	for _, res := range a.Dependencies {
		a.AssignCommonResourceOptions(res)

		if err := res.Create(); err != nil {
			return err
		}
	}
	return nil
}

func (a *Application) DeleteDependencies() error {
	for _, res := range a.Dependencies {
		a.AssignCommonResourceOptions(res)

		if err := res.Delete(); err != nil {
			return err
		}
	}
	return nil
}

func (a *Application) AssignCommonResourceOptions(res *resource.Resource) {
	r := res.Common()
	app := a.Common()

	r.Account = aws.AccountId()
	r.ClusterName = app.ClusterName
	r.Name = app.ServiceAccount
	r.Namespace = app.Namespace
	r.Region = aws.Region()
}
