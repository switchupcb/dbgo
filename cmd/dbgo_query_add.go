package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	subcommand_description_add = "Adds an SQL file _(with the same name as the template file [e.g., `filename.sql`])_ containing an SQL statement _(returned from the `SQL()` function in `filename.go`)_ to the queries directory."
)

// cmdAdd represents the add command
var cmdAdd = &cobra.Command{
	Use:   "add",
	Short: "Add a type-safe SQL file to your queries directory.",
	Long:  subcommand_description_add,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
	},
}

func init() {
	cmdQuery.AddCommand(cmdAdd)
}
