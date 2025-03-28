/*
Copyright © 2025 SwitchUpCB
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/switchupcb/dbgo/cmd/constant"
)

const (
	programDescription = `dbgo generates a database consumer package for your database and domain models (i.e., Go types).`
)

// cmdDBGO represents the base command when called without any subcommands.
var cmdDBGO = &cobra.Command{
	Use:   "dbgo",
	Short: programDescription,
	Long:  programDescription,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   false,
		DisableDescriptions: false,
		HiddenDefaultCmd:    false,
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
//
// called by main.main() once.
func Execute() {
	if err := cmdDBGO.Execute(); err != nil {
		os.Exit(constant.OSExitCodeError)
	}
}

func init() {
	// Persistent Flags work for this command and all subcommands.
	cmdDBGO.PersistentFlags().StringVarP(
		flagYML,
		flagYMLName,
		flagYMLShorthand,
		"",
		flagYMLUsage,
	)
}
