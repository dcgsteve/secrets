package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var deleteProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "delete project",
	Long:  "Deletes a project from the defined Vault instance",
	Args:  cobra.ExactArgs(1),
	Run:   deleteProject,
}

func init() {
	deleteProjectCmd.Flags().Bool("force", false, "force deletion without confirmation")
	deleteCmd.AddCommand(deleteProjectCmd)
}

func deleteProject(cmd *cobra.Command, args []string) {

	checkConfig()

	client, e := getClient()
	if e != nil {
		stopFatal("Failed to create Vault client: ", e.Error())
	}

	// forced?
	f, e := cmd.Flags().GetBool("force")
	if e != nil {
		stopFatal("Error retrieving status of forced flag! Stopping :(")
	}

	// confirm?
	if !f {
		fmt.Printf("Warning: If the project %q exists it will be completely deleted!\nPress Y and ENTER to continue, or anything else to quit:\n", args[0])
		reader := bufio.NewReader(os.Stdin)
		r, _, _ := reader.ReadRune()
		if strings.ToUpper(string(r)) != "Y" {
			stop("Command cancelled - no changes made")
		}

	}

	// retrieve list of keys to delete
	s, e := client.Logical().List(fmt.Sprintf("%s/%s", sc.Store, args[0]))
	if e != nil {
		stopFatal("Could not retrieve list of secrets! Stopping :(")
	}
	if s != nil {
		for _, v := range s.Data["keys"].([]interface{}) {
			_, e = client.Logical().Delete(fmt.Sprintf("%s/%s/%s", sc.Store, args[0], v))
			if e != nil {
				fmt.Printf("Failed to delete secret %q !\n", v)
			}
		}
	}

}
