package cmd

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/util"

	"github.com/spf13/cobra"
)

var utilCmd = &cobra.Command{
	Use:     "utils",
	Short:   "utility commands",
	Aliases: []string{"util"},
}

var enablePrefixAssignmentCmd = &cobra.Command{
	Use:     "enable-prefix-assignment",
	Short:   "Enable Prefix Assignmenr with VPC CNI",
	Aliases: []string{"enable-prefix"},
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster, err := aws.EksDescribeCluster(clusterName)
		if err != nil {
			return err
		}
		return util.EnablePrefixAssignment(cluster)
	},
}

var enableSecurityGroupsForPodsCmd = &cobra.Command{
	Use:     "enable-sg-for-pods",
	Short:   "Enable Security Groups for Pods with VPC CNI",
	Aliases: []string{"enable-sgpods"},
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster, err := aws.EksDescribeCluster(clusterName)
		if err != nil {
			return err
		}
		return util.EnableSecurityGroupsForPods(cluster)
	},
}

var tagSubnetsCmd = &cobra.Command{
	Use:     "tag-subnets",
	Short:   "add kubernetes.io/cluster/<cluster-name> tag to private VPC subnets",
	Aliases: []string{"tag-subnet"},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := aws.EksDescribeCluster(clusterName)
		if err != nil {
			return err
		}

		return util.TagSubnets(clusterName)
	},
}

func init() {
	rootCmd.AddCommand(utilCmd)
	utilCmd.AddCommand(tagSubnetsCmd)
	utilCmd.AddCommand(enablePrefixAssignmentCmd)
	utilCmd.AddCommand(enableSecurityGroupsForPodsCmd)

	tagSubnetsCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "cluster (required)")
	tagSubnetsCmd.MarkFlagRequired("cluster")

	enablePrefixAssignmentCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "cluster (required)")
	enablePrefixAssignmentCmd.MarkFlagRequired("cluster")

	enableSecurityGroupsForPodsCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "cluster (required)")
	enableSecurityGroupsForPodsCmd.MarkFlagRequired("cluster")
}
