package cmd

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/helm"
	"eksdemo/pkg/kubernetes"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource/amg"
	"eksdemo/pkg/resource/amp"
	"eksdemo/pkg/resource/cloudformation"
	"eksdemo/pkg/resource/cluster"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/resource/nodegroup"
	"eksdemo/pkg/resource/organization"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var output printer.Output
var clusterName string

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "view resource(s)",
}

var getHelmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Helm releases",
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster, err := aws.EksDescribeCluster(clusterName)
		if err != nil {
			return err
		}

		kubeContext, err := kubernetes.KubeContextForCluster(cluster)
		if err != nil {
			return err
		}
		if kubeContext == "" {
			return fmt.Errorf("cluster \"%s\" not found in Kubeconfig", clusterName)
		}

		releases, err := helm.List(kubeContext)
		if err != nil {
			return err
		}

		return output.Print(os.Stdout, printer.NewHelmPrinter(releases))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.AddCommand(amg.NewResource().NewGetCmd())
	getCmd.AddCommand(amp.NewResource().NewGetCmd())
	getCmd.AddCommand(cloudformation.NewResource().NewGetCmd())
	getCmd.AddCommand(cluster.NewResource().NewGetCmd())
	getCmd.AddCommand(irsa.NewResource().NewGetCmd())
	getCmd.AddCommand(nodegroup.NewResource().NewGetCmd())
	getCmd.AddCommand(organization.NewResource().NewGetCmd())

	getHelmCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "cluster (required)")
	getHelmCmd.MarkFlagRequired("cluster")
	getHelmCmd.Flags().VarP(cmd.NewOutputFlag(&output), "output", "o", "output format (json|table|yaml)")
	getCmd.AddCommand(getHelmCmd)
}
