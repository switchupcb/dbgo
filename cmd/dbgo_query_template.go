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
	subcommandDescriptionTemplate = "Adds a `name` template to the queries `templates` directory. The template contains Go type database models you can use to return a type-safe SQL statement from the `SQL()` function in `name.go` which is  called by `dbgo query save`."
)

// cmdTemplate represents the dbgo query template command.
var cmdTemplate = &cobra.Command{
	Use:   "template",
	Short: "Add an SQL generator template to your queries directory.",
	Long:  subcommandDescriptionTemplate,
	Run: func(cmd *cobra.Command, args []string) {
		// parse the "--yml" flag.
		yml, err := parseFlagYML()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf("%w", err))

			os.Exit(constant.OSExitCodeError)
		}

		// add each template to the filepath arguments when no filepath arguments are provided.
		if len(args) == 0 {
			queriesGoTemplatesDir := filepath.Join(yml.Generated.Input.Queries, constant.DirnameQueriesTemplates)
			files, err := os.ReadDir(queriesGoTemplatesDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf("%w", err))

				os.Exit(constant.OSExitCodeError)
			}

			for i := range files {
				args = append(args, filepath.Join(queriesGoTemplatesDir, files[i].Name()))
			}
		}

		// run the generator for each filepath argument.
		for _, arg := range args {
			abspath, err := parseArgFilepath(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n\n", fmt.Errorf("%w", err))

				continue
			}

			templateName := filepath.Base(abspath)

			templateDirpath := filepath.Join(
				yml.Generated.Input.Queries,      // queries
				constant.DirnameQueriesTemplates, // templates
				templateName,                     // template (name)
			)

			fmt.Printf("UPDATING TEMPLATE %q at %v\n", templateName, templateDirpath)

			if err := query.Template(templateName, *yml); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n\n", fmt.Errorf("%w", err))

				continue
			}

			fmt.Printf("UPDATED TEMPLATE %q at %v\n\n", templateName, templateDirpath)
		}
	},
}

func init() {
	cmdQuery.AddCommand(cmdTemplate)
}
