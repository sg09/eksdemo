package cmd

import (
	"eksdemo/pkg/eksctl"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

var useContextCmd = &cobra.Command{
	Use:     "use-context CLUSTER",
	Short:   "set currect kubeconfig context",
	Aliases: []string{"uc"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: maybe just context and then check for match and if not write it
		clusterName := args[0]
		eksctlClusterName := eksctl.GetClusterName(clusterName)

		config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			&clientcmd.ConfigOverrides{},
		)

		raw, _ := config.RawConfig()

		for name, context := range raw.Contexts {
			if context.Cluster == eksctlClusterName {
				raw.CurrentContext = name
				// TODO: check error
				clientcmd.ModifyConfig(config.ConfigAccess(), raw, false)
				fmt.Printf("Context switched to: %s\n", name)
				// cluster.Get()
				return nil
			}
		}

		return fmt.Errorf("context not found for cluster: %s", clusterName)
	},
}

func init() {
	rootCmd.AddCommand(useContextCmd)
}
