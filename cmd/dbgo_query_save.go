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
	subcommandDescriptionSave = "Saves an SQL file (with the same name as the template [e.g., `name.sql`]) containing an SQL statement (returned from the `SQL()` function in `name.go`) to the queries directory."
)

// cmdSave represents the dbgo query save command.
var cmdSave = &cobra.Command{
	Use:   "save",
	Short: "Save a type-safe SQL file to your queries directory.",
	Long:  subcommandDescriptionSave,
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

			sqlFilename := templateName + constant.FileExtSQL

			fmt.Printf("SAVING QUERY %v from template at %v\n", sqlFilename, abspath)

			if err := query.Save(templateName, *yml); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n\n", fmt.Errorf("%w", err))

				continue
			}

			fmt.Printf("%v QUERY SAVED from template at %v\n\n", sqlFilename, abspath)
		}
	},
}

func init() {
	cmdQuery.AddCommand(cmdSave)
}
