package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list available secrets",
	Long:  "Produces a list of all secrets (without values) for the defined Vault instance for the defined Project",
	Run:   listSecret,
}

func init() {
	listCmd.Flags().StringP("project", "p", "", "override the currently configured project name")
	rootCmd.AddCommand(listCmd)
}

func listSecret(cmd *cobra.Command, args []string) {

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
	s, e := client.Logical().List(fmt.Sprintf("%s/%s", sc.Store, p))
	if e != nil {
		stop(fmt.Sprintf("Could not check for secrets due to error: %s! (Store was %q, Project was %q)", e.Error(), sc.Store, sc.Project))
	}

	// if we have a result, display list
	if s != nil {
		for _, v := range s.Data["keys"].([]interface{}) {
			fmt.Println(v)
		}
	}

}
