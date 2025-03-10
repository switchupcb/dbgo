package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	query "github.com/switchupcb/dbgo/cmd/dbgo_query"
)

const (
	subcommand_description_gen = "Generates SQL queries for Read (Select) operations and Create (Insert), Update, Delete operations."
)

// cmdQueryGen represents the dbgo query gen command.
var cmdQueryGen = &cobra.Command{
	Use:   "gen",
	Short: "Generates SQL statements from your database.",
	Long:  subcommand_description_gen,
	Run: func(cmd *cobra.Command, args []string) {
		// check for unexpected arguments
		if len(args) != 0 {
			args_string := strings.Join(args, " ")
			fmt.Fprintf(os.Stderr, "Unexpected arguments found: %q", args_string)

			os.Exit(1)
		}

		// parse the "--yml" flag.
		yml, err := parseYML()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf("%w", err))

			os.Exit(1)
		}

		// Run the generator.
		output, err := query.Gen(*yml)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf("%w", err))

			os.Exit(1)
		}

		fmt.Println(output)
	},
}

func init() {
	cmdQuery.AddCommand(cmdQueryGen)
}
