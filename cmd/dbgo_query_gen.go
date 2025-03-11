package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/switchupcb/dbgo/cmd/constant"
	query "github.com/switchupcb/dbgo/cmd/dbgo_query"
)

const (
	subcommandDescriptionGen = "Generates SQL queries for Read (Select) operations and Create (Insert), Update, Delete operations."
)

// cmdQueryGen represents the dbgo query gen command.
var cmdQueryGen = &cobra.Command{
	Use:   "gen",
	Short: "Generates SQL statements from your database.",
	Long:  subcommandDescriptionGen,
	Run: func(cmd *cobra.Command, args []string) {
		// check for unexpected arguments
		if len(args) != 0 {
			argsString := strings.Join(args, " ")

			fmt.Fprintf(os.Stderr, "Unexpected arguments found: %q", argsString)

			os.Exit(constant.OSExitCodeError)
		}

		// parse the "--yml" flag.
		yml, err := parseFlagYML()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf("%w", err))

			os.Exit(constant.OSExitCodeError)
		}

		// Run the generator.
		if err := query.Gen(*yml); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf("%w", err))

			os.Exit(constant.OSExitCodeError)
		}

		fmt.Println("Generated CRUD SQL files at \"" + yml.Generated.Input.Queries + "\"")
	},
}

func init() {
	cmdQuery.AddCommand(cmdQueryGen)
}
