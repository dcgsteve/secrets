package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "retrieve a secret",
	Long:  "Allows for a secret to be read from the defined Vault instance for the defined Project",
	Run:   getSecret,
}

func init() {
	getCmd.Flags().StringP("project", "p", "", "override the currently configured project name")
	rootCmd.AddCommand(getCmd)
}

func getSecret(cmd *cobra.Command, args []string) {

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
	s, _ := client.Logical().Read(fmt.Sprintf("%s/%s/%s", sc.Store, p, args[0]))
	if s == nil {
		stop(fmt.Sprintf("Could not find secret %q! (Store was %q, Project was %q)", args[0], sc.Store, sc.Project))
	}

	// display value (note: we only ever use a single entry here called "value" to simplify usage)
	fmt.Println(s.Data["value"])

}
