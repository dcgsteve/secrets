package cmd

import (
	"fmt"
	"strings"

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
	getCmd.Flags().BoolP("display-key", "k", false, "if the secret is in the format key:value then display the key only")
	getCmd.Flags().BoolP("display-value", "v", false, "if the secret is in the format key:value then display the value only")
	rootCmd.AddCommand(getCmd)
}

func getSecret(cmd *cobra.Command, args []string) {

	checkConfig()

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
		stopFatal("Failed to create Vault client: ", e.Error())
	}

	// get secret map
	s, e := client.Logical().Read(fmt.Sprintf("%s/%s/%s", sc.Store, p, args[0]))
	if e != nil {
		stopFatal(fmt.Sprintf("Could not access Vault correctly! (Store was %q, Project was %q)", sc.Store, sc.Project))
	}
	if s == nil {
		stop(fmt.Sprintf("Could not find secret %q! (Store was %q, Project was %q)", args[0], sc.Store, sc.Project))
	}

	// get secret
	sc := fmt.Sprintf("%v", s.Data["value"])

	// display only part?
	cl := -1
	if b, _ := cmd.Flags().GetBool("display-key"); b {
		cl = 0
	}
	if b, _ := cmd.Flags().GetBool("display-value"); b {
		cl = 1
	}

	// display relevant
	if cl >= 0 {
		s := strings.Split(sc, ":")
		if len(s) != 2 {
			stopFatal("Can only use --display-key or --display-value on key:pair secrets!")
		} else {
			fmt.Println(s[cl])
		}
	} else {
		fmt.Println(sc)
	}

}
