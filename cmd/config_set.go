package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// configSetCmd represents the set command
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "set configuration details",
	Long: `The configuration is encrypted in a local file - this command allows you to
set one or more configuration items by using the appropriate flags`,
	Run: setConfig,
}

func init() {
	configCmd.AddCommand(configSetCmd)

	configSetCmd.Flags().StringP("address", "a", "", "the address of Vault, e.g. http://127.0.0.1:9000")
	configSetCmd.Flags().StringP("token", "t", "", "a valid authorised token for Vault")
	configSetCmd.Flags().StringP("project", "p", "", "a project name (without spaces)")
	configSetCmd.Flags().StringP("store", "s", "", "the Key Value store in Vault to use")
}

func setConfig(cmd *cobra.Command, args []string) {

	a, _ := cmd.Flags().GetString("address")
	if a != "" {
		sc.VaultAddress = a
	}

	t, _ := cmd.Flags().GetString("token")
	if t != "" {
		sc.AuthToken = t
	}

	p, _ := cmd.Flags().GetString("project")
	if p != "" {
		sc.Project = p
	}

	s, _ := cmd.Flags().GetString("store")
	if s != "" {
		sc.Store = s
	}

	e := saveConfig()
	if e != nil {
		log.Fatalf("Failed to save the new configuration: %s", e)
	}

}
