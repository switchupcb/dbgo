/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	subcommand_description_gen = "Generates SQL statements for Read (Select) operations and adds Stored Procedures for Create (Insert), Update, Delete operations to the database."
)

// cmdQueryGen represents the fgen command
var cmdQueryGen = &cobra.Command{
	Use:   "gen",
	Short: "Generates SQL statements from your database.",
	Long:  subcommand_description_gen,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("query gen called")
	},
}

func init() {
	cmdQuery.AddCommand(cmdQueryGen)
}
