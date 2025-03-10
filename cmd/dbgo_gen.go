package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/switchupcb/dbgo/cmd/constant"
	gen "github.com/switchupcb/dbgo/cmd/dbgo_gen"
)

// cmdGen represents the dbgo gen command.
var cmdGen = &cobra.Command{
	Use:   "gen",
	Short: "Use `dbgo gen --yml path/to/yml` to generate a database consumer package.",
	Long:  "Use `dbgo gen --yml path/to/yml` to generate a database consumer package based on domain types, a database, and SQL queries.",
	Run: func(cmd *cobra.Command, args []string) {
		// check for unexpected arguments
		if len(args) != 0 {
			args_string := strings.Join(args, " ")
			fmt.Fprintf(os.Stderr, "Unexpected arguments found: %q", args_string)

			if args_string == cmdQuery.Use {
				fmt.Printf("\n\nDid you mean dbgo %v gen?", cmdQuery.Use)
			}

			os.Exit(constant.OSExitCodeError)
		}

		// parse the "--yml" flag.
		yml, err := parseFlagYML()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf("%w", err))

			os.Exit(constant.OSExitCodeError)
		}

		// Run the generator.
		output, err := gen.Run(*yml)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf("%w", err))

			os.Exit(constant.OSExitCodeError)
		}

		fmt.Println(output)
	},
}

func init() {
	cmdDBGO.AddCommand(cmdGen)
}
