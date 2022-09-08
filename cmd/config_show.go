package cmd

import (
	"fmt"
	"log"

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
	e := getConfig()
	if e != nil {
		log.Fatalf("Unable to read configuration details: %s", e)
	}

	fmt.Printf("Vault address: %s\n", sc.VaultAddress)
	if sc.AuthToken != "" {
		fmt.Println("Vault token: set")
	} else {
		fmt.Println("Vault token: not set")
	}
}
