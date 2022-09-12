package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// correct length of Vault token
const toklen int = 95

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
	printCI("Vault Address", sc.VaultAddress)
	printCI("Vault K/V Store", sc.Store)
	printCI("Current Project", sc.Project)

	// token (redacted)
	if sc.AuthToken != "" {
		l := len(sc.AuthToken)
		if l == toklen {
			fmt.Printf("Vault token (redacted): %s%s%s\n", sc.AuthToken[0:6], strings.Repeat("*", l-12), sc.AuthToken[l-6:l])
		} else {
			fmt.Printf("Vault token: available but looks the wrong length - should be %d characters long ?)\n", toklen)
		}
	} else {
		fmt.Println("Vault token: not entered")
	}

	printCI("Vault Username", sc.Username)

	if sc.Password != "" {
		fmt.Println("Vault password: entered")
	} else {
		fmt.Println("Vault password: not entered")
	}

}

func printCI(k, v string) {
	if v != "" {
		fmt.Printf("%s: %s\n", k, v)
	} else {
		fmt.Printf("%s: not entered\n", k)
	}
}
