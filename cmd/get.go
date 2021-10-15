package cmd

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/helm"
	"eksdemo/pkg/kubernetes"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource/cluster"
	"eksdemo/pkg/resource/irsa"
	"eksdemo/pkg/resource/nodegroup"
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

var getClustersCmd = &cobra.Command{
	Use:     "clusters",
	Short:   "EKS clusters",
	Aliases: []string{"cluster"},
	RunE: func(cmd *cobra.Command, args []string) error {
		clusters, err := cluster.Get()

		if err != nil {
			return err
		}

		currentClusterUrl := kubernetes.GetCurrentContextClusterURL()

		return output.Print(os.Stdout, cluster.NewPrinter(clusters, currentClusterUrl))
	},
}

var getHelmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Helm releases",
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster, err := aws.EksDescribeCluster(clusterName)
		if err != nil {
			return err
		}

		kubeContext, err := kubernetes.GetKubeContextForCluster(cluster)
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

var getIrsaCmd = &cobra.Command{
	Use:   "irsa",
	Short: "IAM Roles for Service Accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster, err := aws.EksDescribeCluster(clusterName)
		if err != nil {
			return err
		}

		roles, err := irsa.Get(cluster)
		if err != nil {
			return err
		}

		return output.Print(os.Stdout, irsa.NewPrinter(roles))
	},
}

var getNodeGroupsCmd = &cobra.Command{
	Use:     "nodegroups",
	Short:   "Managed Node Groups",
	Aliases: []string{"nodegroup", "ng"},
	RunE: func(cmd *cobra.Command, args []string) error {
		nodegroups, err := nodegroup.Get(clusterName)
		if err != nil {
			return err
		}

		return output.Print(os.Stdout, nodegroup.NewPrinter(nodegroups))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getClustersCmd)
	getCmd.AddCommand(getHelmCmd)
	getCmd.AddCommand(getIrsaCmd)
	getCmd.AddCommand(getNodeGroupsCmd)

	getClustersCmd.Flags().VarP(cmd.NewOutputFlag(&output), "output", "o", "output format (json|table|yaml)")

	getHelmCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "cluster (required)")
	getHelmCmd.MarkFlagRequired("cluster")
	getHelmCmd.Flags().VarP(cmd.NewOutputFlag(&output), "output", "o", "output format (json|table|yaml)")

	getIrsaCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "cluster (required)")
	getIrsaCmd.MarkFlagRequired("cluster")
	getIrsaCmd.Flags().VarP(cmd.NewOutputFlag(&output), "output", "o", "output format (json|table|yaml)")

	getNodeGroupsCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "cluster (required)")
	getNodeGroupsCmd.MarkFlagRequired("cluster")
	getNodeGroupsCmd.Flags().VarP(cmd.NewOutputFlag(&output), "output", "o", "output format (json|table|yaml)")
}
