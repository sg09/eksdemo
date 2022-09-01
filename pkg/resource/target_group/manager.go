package target_group

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
	tgOptions, ok := options.(*TargeGroupOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to TargeGroupOptions")
	}

	vpcId := aws.StringValue(tgOptions.Cluster.ResourcesVpcConfig.VpcId)

	if m.DryRun {
		return m.dryRun(tgOptions, vpcId)
	}

	if err := aws.ELBCreateTargetGroup(tgOptions.Name, tgOptions.Protocol, tgOptions.TargetType, vpcId, 1); err != nil {
		return err
	}
	fmt.Printf("Target Group %q created successfully\n", tgOptions.Name)

	return nil
}

func (m *Manager) Delete(options resource.Options) error {
	name := options.Common().Name

	tg, err := m.Getter.GetTargetGroupByName(name)
	if err != nil {
		return err
	}

	if err := aws.ELBDeleteTargetGroup(aws.StringValue(tg.TargetGroupArn)); err != nil {
		return err
	}
	fmt.Printf("Target Group %q deleted\n", name)

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) dryRun(options *TargeGroupOptions, vpcId string) error {
	fmt.Println("\nTarget Group Resource Manager Dry Run:")

	fmt.Printf("Elastic Load Balancing API Call %q with request parameters:\n", "CreateTargetGroup")
	fmt.Printf("Name: %q\n", options.Name)
	fmt.Printf("Port: %q\n", "1")
	fmt.Printf("Protocol: %q\n", options.Protocol)
	fmt.Printf("TargetType: %q\n", options.TargetType)
	fmt.Printf("VpcId: %q\n\n", vpcId)

	return nil
}
