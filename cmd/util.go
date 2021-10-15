package cmd

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/util"

	"github.com/spf13/cobra"
)

var utilCmd = &cobra.Command{
	Use:   "util",
	Short: "utility commands",
}

var tagSubnetsCmd = &cobra.Command{
	Use:     "tag-subnets",
	Short:   "add kubernetes.io/cluster/<cluster-name> tag to private VPC subnets",
	Aliases: []string{"cluster"},
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

	tagSubnetsCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "cluster (required)")
	tagSubnetsCmd.MarkFlagRequired("cluster")
}
