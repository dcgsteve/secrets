package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set configuration details",
	Long: `The configuration is encrypted in a local file - this command allows you to
set one or more configuration items by using the appropriate flags`,
	Run: setConfig,
}

func init() {
	configCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	setCmd.Flags().StringP("address", "a", "", "the address of Vault, e.g. http://127.0.0.1:9000")
	setCmd.Flags().StringP("token", "t", "", "a valid authorised token for Vault")
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

	e := saveConfig()
	if e != nil {
		log.Fatalf("Failed to save the new configuration: %s", e)
	}

}
