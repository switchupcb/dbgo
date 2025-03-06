package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	query "github.com/switchupcb/dbgo/cmd/dbgo_query"
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
		// check for expected arguments.
		if len(args) == 0 {
			fmt.Fprint(os.Stderr, flag_filepath_usage_unspecified)

			os.Exit(1)
		}

		// parse the "-yml" flag.
		yml, err := parseYML(cmdDBGO.PersistentFlags().Lookup(flag_yml_name))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf("%w", err))

			os.Exit(1)
		}

		// run the generator for each filepath argument.
		for _, arg := range args {
			abspath, err := parseArgFilepath(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf("%w", err))

				continue
			}

			output, err := query.Add(abspath, *yml)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf("%w", err))

				continue
			}

			fmt.Println(output)
		}
	},
}

func init() {
	cmdQuery.AddCommand(cmdAdd)
}
