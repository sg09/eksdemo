package nodegroup

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/resource"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type Manager struct {
	Eksctl resource.Manager
	Getter
	resource.EmptyInit
}

func (m *Manager) Create(options resource.Options) error {
	return m.Eksctl.Create(options)
}

func (m *Manager) Delete(options resource.Options) error {
	return m.Eksctl.Delete(options)
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	ngOptions, ok := options.(*NodegroupOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to NodegroupOptions")
	}

	cluster := options.Common().ClusterName
	nodegroup := ngOptions.NodegroupName

	ng, err := m.GetNodeGroupByName(nodegroup, cluster)
	if err != nil {
		return err
	}

	unsetFlags := 0
	update := ""

	if cmd.Flags().Changed("nodes") {
		update += fmt.Sprintf("%d Nodes", ngOptions.UpdateDesired)
	} else {
		ngOptions.UpdateDesired = int(aws.Int64Value(ng.ScalingConfig.DesiredSize))
		unsetFlags++
	}

	if cmd.Flags().Changed("min") {
		if len(update) > 0 {
			update += ", "
		}
		update += fmt.Sprintf("%d Min", ngOptions.MinSize)
	} else {
		ngOptions.UpdateMin = int(aws.Int64Value(ng.ScalingConfig.MinSize))
		unsetFlags++
	}

	if cmd.Flags().Changed("max") {
		if len(update) > 0 {
			update += ", "
		}
		update += fmt.Sprintf("%d Max", ngOptions.MaxSize)
	} else {
		ngOptions.UpdateMax = int(aws.Int64Value(ng.ScalingConfig.MaxSize))
		unsetFlags++
	}

	if unsetFlags == 3 {
		return fmt.Errorf("at least one flag %s is required", strings.Join([]string{"\"nodes\"", "\"min\"", "\"max\""}, ", "))
	}

	fmt.Printf("Updating nodegroup with %s...", update)

	err = aws.EksUpdateNodegroupConfig(cluster, nodegroup, ngOptions.UpdateDesired, ngOptions.UpdateMin, ngOptions.UpdateMax)
	if err != nil {
		return err
	}
	fmt.Println("done")

	return nil
}

func (m *Manager) SetDryRun() {
	m.Eksctl.SetDryRun()
}
