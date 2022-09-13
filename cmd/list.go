package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list secrets or projects",
	Long:  "Produces a list of all secrets or projects for the defined Vault instance for the defined Project",
}

func init() {
	rootCmd.AddCommand(listCmd)
}
