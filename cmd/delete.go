package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "deletes secrets or projects",
	Long:  "Allows for a secret or project to be deleted from the defined Vault instance for the defined Project",
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
