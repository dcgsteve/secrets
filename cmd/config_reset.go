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
	resetCmd.Flags().String("commit", "", "set to YES to allow the reset to run")
}

func resetConfig(cmd *cobra.Command, args []string) {

	commit, e := cmd.Flags().GetString("commit")
	if e != nil {
		log.Fatalf("Invalid commit value: %s", e)
	}

	if commit == "YES" {
		e := setConfigDefaults()
		if e != nil {
			log.Fatalf("Failed to reset default configuration: %s", e)
		}
	} else {
		fmt.Println("Cannot reset configuration as 'commit' flag is not set to YES - no changes have been made")
	}

}
