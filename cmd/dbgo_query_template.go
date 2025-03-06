package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	query "github.com/switchupcb/dbgo/cmd/dbgo_query"
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

			output, err := query.Template(abspath, *yml)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf("%w", err))

				continue
			}

			fmt.Println(output)
		}
	},
}

func init() {
	cmdQuery.AddCommand(templateCmd)
}
