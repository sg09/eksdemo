package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version of eksdemo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version 0.1.0-alpha-2021-10-14")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
