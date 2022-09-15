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
	configSetCmd.Flags().StringP("project", "p", "", "a project name (without spaces)")
	configSetCmd.Flags().StringP("store", "s", "", "the Key Value store in Vault to use")
	configSetCmd.Flags().StringP("username", "u", "", "the Vault username")
	configSetCmd.Flags().StringP("password", "w", "", "the Vault password for the username")

}

func setConfig(cmd *cobra.Command, args []string) {

	var d string
	if d, _ = cmd.Flags().GetString("address"); d != "" {
		sc.VaultAddress = d
	}
	if d, _ = cmd.Flags().GetString("project"); d != "" {
		sc.Project = d
	}
	if d, _ = cmd.Flags().GetString("store"); d != "" {
		sc.Store = d
	}
	if d, _ = cmd.Flags().GetString("username"); d != "" {
		sc.Username = d
	}
	if d, _ = cmd.Flags().GetString("password"); d != "" {
		sc.Password = d
	}

	e := saveConfig()
	if e != nil {
		log.Fatalf("Failed to save the new configuration: %s", e)
	}

}
