package load_balancer

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"fmt"

	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun bool
	Getter
}

func (m *Manager) Create(options resource.Options) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) Delete(options resource.Options) (err error) {
	lbName := options.Common().Name

	elbs, err := m.Getter.GetLoadBalancers(lbName)
	if err != nil {
		return err
	}

	if len(elbs.V1) > 0 {
		err = aws.ELBDeleteLoadBalancerV1(lbName)
	} else {
		err = aws.ELBDeleteLoadBalancerV2(aws.StringValue(elbs.V2[0].LoadBalancerArn))
	}

	if err != nil {
		return err
	}
	fmt.Println("Load balancer deleted successfully")

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}
