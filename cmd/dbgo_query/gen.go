package query

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/switchupcb/dbgo/cmd/config"
)

const (
	generatedQueriesDirname           = "dbgoquerygentemp"
	generatedQueriesOutputDirname     = "output"
	generatedQueriesSchemaSQLFilename = "schema.sql"
)

// Gen runs dbgo query programmatically using the given YML.
func Gen(yml config.YML) (string, error) {
	if yml.Generated.Input.DB.Connection == "" {
		return "", errors.New(err_database_unspecified)
	}

	if yml.Generated.Input.DB.Connection[0] == databaseConnectionEnvironmentVariableSymbol {
		yml.Generated.Input.DB.Connection = os.Getenv(yml.Generated.Input.DB.Connection[1:])
	}

	generatedQueriesFilepath := filepath.Join(
		yml.Generated.Input.Queries, // queries
		generatedQueriesDirname,     // generatedQueriesDirname
	)

	if _, err := os.Stat(generatedQueriesFilepath); err == nil {
		return "", fmt.Errorf("warning: directory at must be deleted: %q", generatedQueriesFilepath)
	} else if errors.Is(err, os.ErrNotExist) {

	} else {
		return "", fmt.Errorf("error checking for directory space: %w", err)
	}

	// Create an sqlc CRUD Generator project.
	if err := os.MkdirAll(generatedQueriesFilepath, writeFileMode); err != nil {
		return "", fmt.Errorf("mkdir all: %w", err)
	}

	// Add schema file.
	pgdump := exec.Command("pg_dump", //nolint:gosec // disable G204
		yml.Generated.Input.DB.Connection,
		"--schema-only",
		"-f", filepath.Join(generatedQueriesFilepath, generatedQueriesSchemaSQLFilename),
	)

	std, err := pgdump.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("write schema file: pg_dump: %q: %w", string(std), err)
	}

	// Add static files.
	if err := os.WriteFile(filepath.Join(generatedQueriesFilepath, file_name_dummy_sql), []byte(file_content_dummy_sql), writeFileMode); err != nil {
		return "", fmt.Errorf("write dummy file: %w", err)
	}

	file_path_sqlc_json := filepath.Join(generatedQueriesFilepath, file_name_sqlc_json)
	if err := os.WriteFile(file_path_sqlc_json, []byte(file_content_sqlc_json), writeFileMode); err != nil {
		return "", fmt.Errorf("write config file: %w", err)
	}

	// Run the CRUD Generator.
	sqlc := exec.Command("sqlc", "generate", "-f", file_path_sqlc_json)

	std, err = sqlc.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("write CRUD SQL: sqlc generate: %q: %w", string(std), err)
	}

	// Output the CRUD SQL to the queries directory.
	src, err := os.ReadFile(filepath.Join(generatedQueriesFilepath, generatedQueriesOutputDirname, file_name_dummy_sql))
	if err != nil {
		return "", fmt.Errorf("read CRUD SQL: %w", err)
	}

	srcQueries := bytes.Split(src, []byte{newline, newline})
	for i := range srcQueries {
		query := srcQueries[i]

		if len(query) == 0 {
			continue
		}

		// name represents the query name (e.g., `InsertUser` in `-- name: InsertUser :one`)
		var name []byte

		colon_count := 0

	parseName:
		for i := range query {
			switch colon_count {
			case 0:
				if query[i] == colon {
					colon_count++
				}
			case 1:
				switch query[i] {
				case whitespace:
				case colon:
					colon_count++
				default:
					name = append(name, query[i])
				}
			case 2: //nolint:mnd
				break parseName
			}
		}

		if colon_count != 2 { //nolint:mnd
			return "", fmt.Errorf("encountered invalid CRUD SQL at statement %d\n%q", i, string(srcQueries[i]))
		}

		if err := os.WriteFile(filepath.Join(yml.Generated.Input.Queries, string(name)+fileExtSQL), query, writeFileMode); err != nil {
			return "", fmt.Errorf("write CRUD SQL FILE at statement %d\n%q", i, string(query))
		}
	}

	if err := os.RemoveAll(generatedQueriesFilepath); err != nil {
		return "", fmt.Errorf("clean: %w", err)
	}

	return fmt.Sprintf("Generated CRUD SQL files at %q", yml.Generated.Input.Queries), nil
}
