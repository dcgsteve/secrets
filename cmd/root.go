package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

type secretsConfig struct {
	VaultAddress string
	AuthToken    string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Team based secret manager for Hashicorp Vault",
	Long:  "Allows easy use in automated processes to store/retrieve simple key/value pairs of information",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	// ensure base minimum config file is available
	if fileNotExists(getConfigFileName()) {
		_, e := setConfigDefaults()
		if e != nil {
			log.Fatalf("No configuration file found and failed to set defaults: %s", e)
		}
	}

}
