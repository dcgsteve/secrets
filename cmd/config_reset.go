package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "reset configuration to defaults",
	Long: `Forces a reset of the configuration details held to the defaults

As a protection the 'commit' flag must be set to YES in order to run the reset

Example:
  secrets reset --commit YES
`,
	Run: resetConfig,
}

func init() {
	configCmd.AddCommand(resetCmd)
	resetCmd.Flags().Bool("confirm", false, "specify this flag to confirm the reset can run")
}

func resetConfig(cmd *cobra.Command, args []string) {

	checkConfig()

	confirm, _ := cmd.Flags().GetBool("confirm")
	if confirm {
		e := setConfigDefaults()
		if e != nil {
			log.Fatalf("Failed to reset default configuration: %s", e)
		}
	} else {
		fmt.Println("Cannot reset configuration as 'confirm' flag was not specified - no changes have been made")
	}

}
