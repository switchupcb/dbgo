package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/switchupcb/dbgo/cmd/constant"
	query "github.com/switchupcb/dbgo/cmd/dbgo_query"
)

const (
	subcommand_description_template = "Adds a `name` template to the queries `templates` directory. The template contains Go type database models you can use to return a type-safe SQL statement from the `SQL()` function in `name.go` which is  called by `db query save`."
)

// cmdTemplate represents the dbgo query template command.
var cmdTemplate = &cobra.Command{
	Use:   "template",
	Short: "Add an SQL generator template to your queries directory.",
	Long:  subcommand_description_template,
	Run: func(cmd *cobra.Command, args []string) {
		// parse the "--yml" flag.
		yml, err := parseFlagYML()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf("%w", err))

			os.Exit(constant.OSExitCodeError)
		}

		// run the generator for each template.
		if len(args) == 0 {
			queriesGoTemplatesDir := filepath.Join(yml.Generated.Input.Queries, constant.DirnameQueriesTemplates)
			files, err := os.ReadDir(queriesGoTemplatesDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf("%w", err))

				os.Exit(constant.OSExitCodeError)
			}

			for i := range files {
				output, err := query.Template(filepath.Join(queriesGoTemplatesDir, files[i].Name()), *yml)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n\n", fmt.Errorf("%w", err))

					continue
				}

				fmt.Printf("%v\n\n", output)
			}

			return
		}

		// run the generator for each filepath argument.
		for _, arg := range args {
			abspath, err := parseArgFilepath(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n\n", fmt.Errorf("%w", err))

				continue
			}

			output, err := query.Template(abspath, *yml)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n\n", fmt.Errorf("%w", err))

				continue
			}

			fmt.Printf("%v\n\n", output)
		}
	},
}

func init() {
	cmdQuery.AddCommand(cmdTemplate)
}
