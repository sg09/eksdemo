package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List EKS clusters",
	Aliases: []string{"ls"},
	Hidden:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return getClustersCmd.RunE(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
