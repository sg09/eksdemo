package application

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/irsa"
)

type IamPolicy struct {
	irsa.PolicyType
	Policy []string
}

func (p *IamPolicy) NewIrsa(options Options) *resource.Resource {
	irsaOptions := &irsa.IrsaOptions{
		CommonOptions: resource.CommonOptions{
			Account:     aws.AccountId(),
			ClusterName: options.Common().ClusterName,
			Name:        options.Common().ServiceAccount,
			Namespace:   options.Common().Namespace,
			Region:      aws.Region(),
		},
		PolicyType: p.PolicyType,
		Policy:     p.Policy,
	}

	res := irsa.NewResource()
	res.Options = irsaOptions

	return res
}

func (p *IamPolicy) Create(options Options) error {
	res := p.NewIrsa(options)

	return res.Create(res.Options)
}

func (p *IamPolicy) Delete(options Options) error {
	res := p.NewIrsa(options)

	return res.Delete(res.Options)
}
