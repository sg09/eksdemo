package cmd

import (
	"eksdemo/cmd/create"
	"eksdemo/cmd/install"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:              "eksdemo",
	Short:            "Create and manage demo clusters",
	PersistentPreRun: preRun,
	SilenceErrors:    true,
	Long: `An opinioned toolkit to quickly and easily create and manage
EKS clusters. Install applications along with required IAM Roles
using best practices configurations.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func preRun(cmd *cobra.Command, args []string) {
	// This will work in the future if the issue below is fixed:
	// https://github.com/spf13/cobra/issues/1413
	// cmd.SilenceUsage = true
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(
		create.NewCreateCmd(),
		newCmdDelete(),
		install.NewInstallCmd(),
		install.NewUninstallCmd(),
		newCmdUpdate(),
	)

	// TODO: implement configuration
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.eksdemo.yaml)")

	// Hide help command
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".eksdemo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".eksdemo")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
