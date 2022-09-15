package cmd

import (
	"fmt"

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
	setCmd.Flags().BoolP("force", "f", false, "allow setting a secret if it already exists")
	rootCmd.AddCommand(setCmd)
}

func setSecret(cmd *cobra.Command, args []string) {

	checkConfig()

	// check syntax
	if len(args) != 2 {
		stop("This command needs two arguments, the first being the name of the secret to write and the second being the value of the secret")
	}

	// parameter takes precendence over current config
	p, _ := cmd.Flags().GetString("project")
	if p == "" {
		p = sc.Project
	}

	// get client
	client := cliGetClient()

	// init data
	d := map[string]interface{}{
		"value": args[1],
	}

	// is secret there?
	f, _ := cmd.Flags().GetBool("force")
	s, e := client.Logical().Read(fmt.Sprintf("%s/%s/%s", sc.Store, p, args[0]))
	if e != nil {
		stopFatal(fmt.Sprintf("Failed to check for secret %q: %s", args[0], e))
	}
	if s != nil && !f {
		stopFatal("Secret exists - cannot overwrite without --force flag")
	}

	// write
	_, e = client.Logical().Write(fmt.Sprintf("%s/%s/%s", sc.Store, p, args[0]), d)
	if e != nil {
		stopFatal("Failed to write secret: ", e.Error())
	}

}
