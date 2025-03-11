package query

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/switchupcb/dbgo/cmd/config"
	"github.com/switchupcb/dbgo/cmd/constant"
)

const (
	interpretedFileContentStatic = "package " + constant.PkgNameSchemaGo +
		`

import . "github.com/switchupcb/jet/v2/postgres"

// SQL returns return an SQL statement.
//
// You can use Jet to write type-safe SQL queries.
//
// Read https://github.com/go-jet/jet#lets-write-some-sql-queries-in-go for more information.
func SQL() (string, error) {
	stmt := *new(SelectStatement)

	query, _ := stmt.Sql()
	
	return query, nil
}
	`
)

// Template runs dbgo query template programmatically using the given template name and YML.
func Template(name string, yml config.YML) error {
	// Copy the existing schema.go file for Go type autocompletion within the template.
	copied := true

	queriesSchemaGoFilepath := filepath.Join(
		yml.Generated.Input.Queries,       // queries
		constant.DirnameQueriesSchema,     // schema
		constant.FilenameTemplateSchemaGo, // schema.go
	)

	templateSchemaGoFilepath := filepath.Join(
		yml.Generated.Input.Queries,       // queries
		constant.DirnameQueriesTemplates,  // templates
		name,                              // template (name)
		constant.FilenameTemplateSchemaGo, // schema.go
	)

	if _, err := os.Stat(queriesSchemaGoFilepath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println(
				"WARNING: The template's schema.go file was not updated because" +
					"schema.go was not found at " + queriesSchemaGoFilepath +
					"\n\tYou can use `dbgo query schema` to regenerate a schema.go file." +
					"\n\tThen, run `dbgo query template` again to update the template's schema.go file.",
			)

			copied = false
		} else {
			return fmt.Errorf("error checking .../queries/schema/schema.go file space: %w", err)
		}
	}

	if copied {
		if err := constant.CopyFile(queriesSchemaGoFilepath, templateSchemaGoFilepath); err != nil {
			return fmt.Errorf("error copying queries schema.go to template: %w", err)
		}
	}

	// Create the interpreted function file when it doesn't exist.
	templateInterpretedGoFilepath := filepath.Join(
		yml.Generated.Input.Queries,      // queries
		constant.DirnameQueriesTemplates, // templates
		name,                             // template (name)
		name+constant.FileExtGo,          // template.go
	)

	if _, err := os.Stat(templateInterpretedGoFilepath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fileContent := []byte(interpretedFileContentStatic)

			if err := os.WriteFile(templateInterpretedGoFilepath, fileContent, constant.FileModeWrite); err != nil {
				return fmt.Errorf("template: write: %w", err)
			}
		} else {
			return fmt.Errorf("error checking template file space: %w", err)
		}
	}

	return nil
}
