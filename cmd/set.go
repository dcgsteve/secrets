package cmd

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "sets a secret",
	Long: `Allows for a secret to be writen to the defined Vault instance for the defined Project

For example: secrets set colour yellow
`,
	Run: setSecret,
}

func init() {
	setCmd.Flags().StringP("project", "p", "", "override the currently configured project name")
	rootCmd.AddCommand(setCmd)
}

func setSecret(cmd *cobra.Command, args []string) {

	// check syntax
	if len(args) != 2 {
		stop("This command needs two arguments, the first being the name of the secret to write and the second being the value of the secret")
	}

	// parameter takes precendence over current config
	p, _ := cmd.Flags().GetString("project")
	if p == "" {
		p = sc.Project
	}

	// init client
	client, e := api.NewClient(&api.Config{Address: sc.VaultAddress, HttpClient: httpClient})
	if e != nil {
		stop(fmt.Sprintf("Failed to create Vault client: %s", e))
	}
	client.SetToken(sc.AuthToken)

	// init data
	d := map[string]interface{}{
		"value": args[1],
	}

	// write
	_, e = client.Logical().Write(fmt.Sprintf("%s/%s/%s", sc.Store, p, args[0]), d)
	if e != nil {
		stop(fmt.Sprintf("Failed to write secret: %s", e))
	}

}
