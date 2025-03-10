package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/switchupcb/dbgo/cmd/constant"
	query "github.com/switchupcb/dbgo/cmd/dbgo_query"
)

const (
	subcommand_description_schema = "Generates a schema.sql and schema.go file representing your database in the queries directory."
)

var (
	cmdSchemaFlagSQL = new(bool)
	cmdSchemaFlagGo  = new(bool)
)

// cmdSchema represents the dbgo query schema command.
var cmdSchema = &cobra.Command{
	Use:   "schema",
	Short: "Generate a schema.sql and schema.go file in your queries directory.",
	Long:  subcommand_description_schema,
	Run: func(cmd *cobra.Command, args []string) {
		// parse the "--yml" flag.
		yml, err := parseFlagYML()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf("%w", err))

			os.Exit(constant.OSExitCodeError)
		}

		// generate the schema files.
		if *cmdSchemaFlagSQL != *cmdSchemaFlagGo {
			if *cmdSchemaFlagSQL {

			}

			if *cmdSchemaFlagGo {

			}
		}

		output, err := query.Schema(*yml)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n\n", fmt.Errorf("%w", err))
		}

		fmt.Printf("%v\n\n", output)
	},
}

func init() {
	cmdQuery.AddCommand(cmdSchema)

	cmdSchemaFlagSQL = cmdSchema.Flags().BoolP("sql", "s", false, "Use --sql to only generate a schema.sql file.")
	cmdSchemaFlagGo = cmdSchema.Flags().BoolP("go", "g", false, "Use --go to only generate a schema.go file.")
}
