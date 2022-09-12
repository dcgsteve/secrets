package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a secret",
	Long:  "Allows for a secret to be deleted from the defined Vault instance for the defined Project",
	Run:   deleteSecret,
}

func init() {
	deleteCmd.Flags().StringP("project", "p", "", "override the currently configured project name")
	rootCmd.AddCommand(deleteCmd)
}

func deleteSecret(cmd *cobra.Command, args []string) {

	// check secret name supplied
	if len(args) == 0 {
		stop("The name of the secret to read was not supplied!")
	}

	// parameter takes precendence over current config
	p, _ := cmd.Flags().GetString("project")
	if p == "" {
		p = sc.Project
	}

	client, e := getClient()
	if e != nil {
		stop("Failed to create Vault client: ", e.Error())
	}

	// get secret map
	_, e = client.Logical().Delete(fmt.Sprintf("%s/%s/%s", sc.Store, p, args[0]))
	if e != nil {
		stop(fmt.Sprintf("Could not access Vault correctly! (Store was %q, Project was %q)", sc.Store, sc.Project))
	}

	fmt.Printf("Secret %q was deleted (if it existed)", args[0])

}
