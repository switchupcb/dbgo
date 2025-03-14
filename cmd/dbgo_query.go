package cmd

import (
	"github.com/spf13/cobra"
)

// cmdQuery represents the dbgo query command.
var cmdQuery = &cobra.Command{
	Use:   "query",
	Short: "Use the `dbgo query` manager to manage SQL statements.",
	Long:  "Use the `dbgo query` manager to save customized type-safe SQL statements or generate them.",
}

func init() {
	cmdDBGO.AddCommand(cmdQuery)
}
