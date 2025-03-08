package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	query "github.com/switchupcb/dbgo/cmd/dbgo_query"
)

const (
	subcommand_description_save = "Saves an SQL file (with the same name as the template [e.g., `name.sql`]) containing an SQL statement (returned from the `SQL()` function in `name.go`) to the queries directory."
)

// cmdSave represents the dbgo query save command.
var cmdSave = &cobra.Command{
	Use:   "save",
	Short: "Save a type-safe SQL file to your queries directory.",
	Long:  subcommand_description_save,
	Run: func(cmd *cobra.Command, args []string) {
		// parse the "--yml" flag.
		yml, err := parseYML()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf("%w", err))

			os.Exit(1)
		}

		// run the generator for each template.
		if len(args) == 0 {
			queriesGoTemplatesDir := filepath.Join(yml.Generated.Input.Queries, queriesGoTemplatesDirname)
			files, err := os.ReadDir(queriesGoTemplatesDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf("%w", err))

				os.Exit(1)
			}

			for i := range files {
				output, err := query.Save(filepath.Join(queriesGoTemplatesDir, files[i].Name()), *yml)
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

			output, err := query.Save(abspath, *yml)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n\n", fmt.Errorf("%w", err))

				continue
			}

			fmt.Printf("%v\n\n", output)
		}
	},
}

func init() {
	cmdQuery.AddCommand(cmdSave)
}
