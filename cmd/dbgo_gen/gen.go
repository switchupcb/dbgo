package gen

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/switchupcb/dbgo/cmd/config"
	"github.com/switchupcb/dbgo/cmd/constant"
)

var (
	newlineBuffer = []byte{constant.Newline}

	// sqlcQueryAnnotationNameIdentifier represents the sqlc Query Annotation name format.
	//
	// https://docs.sqlc.dev/en/stable/reference/query-annotations.html
	sqlcQueryAnnotationNameIdentifier = []byte("-- name: ")

	// sqlcQueryAnnotationCommandExec represents the :exec sqlc query annotation command.
	sqlcQueryAnnotationCommandExec = []byte("exec")
)

// Gen runs dbgo gen programmatically using the given YML.
//
// sqlc is used to generate code from SQL statements.
func Gen(yml config.YML, keepcombined bool) error {
	fmt.Println("")

	// Combine the SQL statements in the queries directory to a single SQL file.
	printQueryAnnotationGuide := false

	files, err := os.ReadDir(yml.Generated.Input.Queries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf("%w", err))

		os.Exit(constant.OSExitCodeError)
	}

	var combinedFileContentSQL []byte

	for i := range files {
		if files[i].IsDir() {
			continue
		}

		if files[i].Name() == constant.FilenameQueriesCombinedSQLKept {
			continue
		}

		if path.Ext(files[i].Name()) == constant.FileExtSQL {
			src, err := os.ReadFile(filepath.Join(yml.Generated.Input.Queries, files[i].Name()))
			if err != nil {
				return fmt.Errorf("error reading sql file: %w", err)
			}

			// check for the sqlc query annotation line by line.
			queryAnnotationExists := false

			srcLines := bytes.Split(src, newlineBuffer)
			for i := range srcLines {
				if len(srcLines[i]) == 0 {
					continue
				}

				if bytes.Contains(srcLines[i], sqlcQueryAnnotationNameIdentifier) {
					queryAnnotationExists = true

					break
				}
			}

			// CURRENT: Add `exec` query annotation to the SQL statement when it doesn't exist.
			//
			// FUTURE: Use a customizable algorithm to automatically add directives when they don't exist.
			if !queryAnnotationExists {
				printQueryAnnotationGuide = true

				fmt.Printf("WARNING: %v does not have a valid query annotation.\n\tUsing :exec by default.\n", files[i].Name())

				queryAnnotationName := []byte(files[i].Name()[:len(files[i].Name())-len(constant.FileExtSQL)])

				// -- name: name :exec
				queryAnnotation := make(
					[]byte,
					0,
					len(sqlcQueryAnnotationNameIdentifier)+
						len(queryAnnotationName)+
						2+ // whitespace + colon
						len(sqlcQueryAnnotationCommandExec),
				)

				queryAnnotation = append(queryAnnotation, sqlcQueryAnnotationNameIdentifier...) // `--name: `
				queryAnnotation = append(queryAnnotation, queryAnnotationName...)               // name
				queryAnnotation = append(queryAnnotation,
					constant.Whitespace, // last ' ' in `--name: name `
					constant.Colon,      // :
				)
				queryAnnotation = append(queryAnnotation, sqlcQueryAnnotationCommandExec...) // exec

				combinedFileContentSQL = append(combinedFileContentSQL, queryAnnotation...)
				combinedFileContentSQL = append(combinedFileContentSQL, constant.Newline)
			}

			// Add the SQL statement to the combined file.
			combinedFileContentSQL = append(combinedFileContentSQL, src...)
			combinedFileContentSQL = append(combinedFileContentSQL, constant.Newline, constant.Newline)
		}
	}

	// Create the sqlc generate project.
	//
	// dirpathGenerationSpace represents the directory used to contain files during generation.
	//
	// The dirpathGenerationSpace directory is used when the sqlc CRUD Generator project is created.
	dirpathGenerationSpace := filepath.Join(
		yml.Generated.Input.Queries, // queries
		constant.DirnameTempQueriesGenerationSQLC,
	)

	// do not overwrite an existing directory from the user.
	if _, err := os.Stat(dirpathGenerationSpace); err == nil {
		return fmt.Errorf("directory at must be deleted: %q", dirpathGenerationSpace)
	} else if errors.Is(err, os.ErrNotExist) {

	} else {
		return fmt.Errorf("error checking directory space: %w", err)
	}

	if err := os.MkdirAll(dirpathGenerationSpace, constant.FileModeWrite); err != nil {
		return fmt.Errorf("mkdir all: %w", err)
	}

	// combined.sql
	filepathCombinedSQL := filepath.Join(
		dirpathGenerationSpace,
		constant.FilenameQueriesCombinedSQL,
	)

	if err := os.WriteFile(
		filepathCombinedSQL,
		combinedFileContentSQL,
		constant.FileModeWrite,
	); err != nil {
		return fmt.Errorf("write combined.sql: %w", err)
	}

	// sqlc.yaml
	srcYML, err := fileContentSQLCYML(yml)
	if err != nil {
		return fmt.Errorf("write config file: relative pathfinder: %w", err)
	}

	filepathSQLCYML := filepath.Join(dirpathGenerationSpace, constant.FilenameSQLConfig)
	if err := os.WriteFile(
		filepathSQLCYML,
		[]byte(srcYML),
		constant.FileModeWrite,
	); err != nil {
		return fmt.Errorf("write sqlc.yaml: %w", err)
	}

	// Run sqlc generate.
	sqlc := exec.Command("sqlc", "generate", "-f", filepathSQLCYML)
	std, err := sqlc.CombinedOutput()
	if err != nil { //nolint:wsl
		return fmt.Errorf("write Go code: sqlc generate: %q: %w", string(std), err)
	}

	// Clean the project.
	if keepcombined {
		filepathKeptCombinedSQL := filepath.Join(
			yml.Generated.Input.Queries,
			constant.FilenameQueriesCombinedSQLKept,
		)

		if err := constant.CopyFile(filepathCombinedSQL, filepathKeptCombinedSQL); err != nil {
			return fmt.Errorf("error copying combined.sql to queries directory: %w", err)
		}
	}

	if err := os.RemoveAll(dirpathGenerationSpace); err != nil {
		return fmt.Errorf("clean: %w", err)
	}

	if printQueryAnnotationGuide {
		fmt.Println("\nTIP: You can read more about sqlc query annotations at https://docs.sqlc.dev/en/stable/reference/query-annotations.html")
	}

	return nil
}
