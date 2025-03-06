package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/switchupcb/dbgo/cmd/gen"
)

var ymlFlag *string

// cmdGen represents the gen command
var cmdGen = &cobra.Command{
	Use:   "gen",
	Short: "Use `dbgo gen -yml path/to/yml` to generate a database consumer package.",
	Long:  "Use `dbgo gen -yml path/to/yml` to generate a database consumer package based on domain types, a database, and SQL queries.",
	Run: func(cmd *cobra.Command, args []string) {
		// todo: https://github.com/spf13/cobra/issues/2250
		if *ymlFlag == "" {
			fmt.Println(flag_yml_usage_unspecified)

			os.Exit(1)
		}

		output, err := gen.Run(*ymlFlag)
		if err != nil {
			fmt.Printf("error: %v", err)

			os.Exit(1)
		}

		fmt.Println(output)
	},
}

func init() {
	cmdDBGO.AddCommand(cmdGen)

	ymlFlag = cmdGen.Flags().String("yml", "", flag_yml_usage)
}
