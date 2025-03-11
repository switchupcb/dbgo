package query

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/switchupcb/dbgo/cmd/config"
	"github.com/switchupcb/dbgo/cmd/constant"
)

// Gen runs dbgo query gen programmatically using the given YML.
func Gen(yml config.YML) error {
	filepathSchemaSQL := filepath.Join(
		yml.Generated.Input.Queries,       // queries
		constant.DirnameQueriesSchema,     // schema
		constant.FilenameQueriesSchemaSQL, // schema.sql
	)

	if _, err := os.Stat(filepathSchemaSQL); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf( //nolint:stylecheck // ST1005
				`schema.sql was not found at %q
	You can use 'db query schema' to regenerate a schema.sql file.
	Then, run 'db query gen' again to generate CRUD SQL.`,
				filepathSchemaSQL,
			)
		} else {
			return fmt.Errorf("error checking .../queries/schema/schema.go file space: %w", err)
		}
	}

	// dirpathGenerationSpace represents the directory used to contain files during generation.
	//
	// The dirpathGenerationSpace directory is used when the sqlc CRUD Generator project is created.
	dirpathGenerationSpace := filepath.Join(
		yml.Generated.Input.Queries, // queries
		constant.DirnameTempQueriesGenerationSQL,
	)

	// Do not overwrite an existing directory from the user.
	if _, err := os.Stat(dirpathGenerationSpace); err == nil {
		return fmt.Errorf("directory at must be deleted: %q", dirpathGenerationSpace)
	} else if errors.Is(err, os.ErrNotExist) {

	} else {
		return fmt.Errorf("error checking directory space: %w", err)
	}

	// Create an sqlc CRUD Generator project directory.
	if err := os.MkdirAll(dirpathGenerationSpace, constant.FileModeWrite); err != nil {
		return fmt.Errorf("mkdir all: %w", err)
	}

	// Add static files.
	//
	// placeholder.sql
	if err := os.WriteFile(
		filepath.Join(
			dirpathGenerationSpace,
			filenamePlaceholderSQL,
		),
		[]byte(fileContentPlaceholderSQL),
		constant.FileModeWrite,
	); err != nil {
		return fmt.Errorf("write placeholder file: %w", err)
	}

	// crud.json
	relativeSchemaPath, err := filepath.Rel(dirpathGenerationSpace, filepathSchemaSQL)
	if err != nil {
		return fmt.Errorf("write config file: relative pathfinder: %w", err)
	}

	filepathSQLCJSON := filepath.Join(
		dirpathGenerationSpace,
		filenameSQLCJSON,
	)

	if err := os.WriteFile(
		filepathSQLCJSON,
		[]byte(fileContentSQLCJSON(relativeSchemaPath)),
		constant.FileModeWrite,
	); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	// Run the CRUD Generator.
	sqlc := exec.Command("sqlc", "generate", "-f", filepathSQLCJSON)
	std, err := sqlc.CombinedOutput()
	if err != nil { //nolint:wsl
		return fmt.Errorf("write CRUD SQL: sqlc generate: %q: %w", string(std), err)
	}

	// Output the CRUD SQL to the queries directory.
	src, err := os.ReadFile(
		filepath.Join(
			dirpathGenerationSpace,
			dirnameGeneratedQueriesOutput,
			filenameGeneratedQueriesSQL,
		),
	)
	if err != nil {
		return fmt.Errorf("read CRUD SQL: %w", err)
	}

	srcQueries := bytes.Split(src, []byte{constant.Newline, constant.Newline})
	for i := range srcQueries {
		query := srcQueries[i]

		if len(query) == 0 {
			continue
		}

		// name represents the query name (e.g., `InsertUser` in `-- name: InsertUser :one`)
		var name []byte

		colonCount := 0

	parseName:
		for i := range query {
			switch colonCount {
			case 0:
				if query[i] == constant.Colon {
					colonCount++
				}
			case 1:
				switch query[i] {
				case constant.Whitespace:
				case constant.Colon:
					colonCount++
				default:
					name = append(name, query[i])
				}
			case 2: //nolint:mnd
				break parseName
			}
		}

		if colonCount != 2 { //nolint:mnd
			return fmt.Errorf("encountered invalid CRUD SQL at statement %d\n%q", i, string(srcQueries[i]))
		}

		if err := os.WriteFile(filepath.Join(yml.Generated.Input.Queries, string(name)+constant.FileExtSQL), query, constant.FileModeWrite); err != nil {
			return fmt.Errorf("write CRUD SQL FILE at statement %d\n%q", i, string(query))
		}
	}

	if err := os.RemoveAll(dirpathGenerationSpace); err != nil {
		return fmt.Errorf("clean: %w", err)
	}

	return nil
}
