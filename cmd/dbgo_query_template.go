package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	subcommand_description_template = "Adds a `filename.go` template file to the queries directory containing database models you can use to return a type-safe SQL statement from the `SQL()` function called by `db query add`."
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Add an SQL generator template file to your queries directory.",
	Long:  subcommand_description_template,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("template called")
	},
}

func init() {
	cmdQuery.AddCommand(templateCmd)
}
