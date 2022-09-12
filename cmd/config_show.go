package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show configuration details",
	Long:  "Displays the SECRETS configuration information held locally",
	Run:   showConfig,
}

func init() {
	configCmd.AddCommand(showCmd)
}

func showConfig(cmd *cobra.Command, args []string) {

	// display
	printCI("Vault Address", sc.VaultAddress, false)
	printCI("Vault K/V Store", sc.Store, false)
	printCI("Current Project", sc.Project, false)
	printCI("Vault Username", sc.Username, false)
	printCI("Vault Password", sc.Password, true)

}

func printCI(k, v string, r bool) {
	if v != "" {
		if r {
			fmt.Printf("%s: [entered]\n", k)
		} else {
			fmt.Printf("%s: %s\n", k, v)
		}

	} else {
		fmt.Printf("%s: [not entered]\n", k)
	}
}
