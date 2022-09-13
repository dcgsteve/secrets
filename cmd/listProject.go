package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var listProjectCmd = &cobra.Command{
	Use:   "projects",
	Short: "list available projects",
	Long:  "Produces a list of all projects for the defined Vault instance for the defined Project",
	Run:   listProject,
}

func init() {
	listCmd.AddCommand(listProjectCmd)
}

func listProject(cmd *cobra.Command, args []string) {

	client, e := getClient()
	if e != nil {
		stop("Failed to create Vault client: ", e.Error())
	}

	// get secret map
	s, e := client.Logical().List(sc.Store)
	if e != nil {
		stop(fmt.Sprintf("Could not check for projects due to error: %s! (Store was %q)", e.Error(), sc.Store))
	}

	// if we have a result, display list
	if s != nil {
		for _, v := range s.Data["keys"].([]interface{}) {
			fmt.Println(v)
		}
	}

}
