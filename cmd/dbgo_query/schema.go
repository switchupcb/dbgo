package query

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/switchupcb/dbgo/cmd/config"
	"github.com/switchupcb/dbgo/cmd/constant"
	"github.com/switchupcb/jet/v2/generator/postgres"
)

// Schema runs dbgo query schema programmatically using the given YML.
func Schema(yml config.YML, schemago, schemasql bool) error {
	var err error

	yml.Generated.Input.DB.Connection, err = validatedDatabaseConnection(yml)
	if err != nil {
		return err
	}

	if yml.Generated.Input.DB.Schema == "" {
		yml.Generated.Input.DB.Schema = constant.DatabaseSchemaNameDefault
	}

	queriesSchemaDir := filepath.Join(
		yml.Generated.Input.Queries,   // queries
		constant.DirnameQueriesSchema, // schema
	)

	if schemago {
		fmt.Printf("\nGenerating schema.go in %q\n", queriesSchemaDir)

		if err := SchemaGo(
			queriesSchemaDir,
			yml.Generated.Input.DB.Connection,
			yml.Generated.Input.DB.Schema,
		); err != nil {
			return fmt.Errorf("error generating schema.go: %w", err)
		}

		fmt.Printf("Generated schema.go in %q\n", queriesSchemaDir)
	}

	if schemasql {
		fmt.Printf("\nGenerating schema.sql in %q\n", queriesSchemaDir)

		if err := SchemaSQL(
			queriesSchemaDir,
			yml.Generated.Input.DB.Connection,
			yml.Generated.Input.DB.Schema,
		); err != nil {
			return fmt.Errorf("error generating schema.go: %w", err)
		}

		fmt.Printf("Generated schema.sql in %q\n", queriesSchemaDir)
	}

	return nil
}

// SchemaSQL generates a schema.sql file in the given directory using the
// database connection string and database schema name.
func SchemaSQL(dirpath, dbconnection, dbschema string) error {
	pgdump := exec.Command("pg_dump", //nolint:gosec // disable G204
		dbconnection,
		"-n", dbschema,
		"--schema-only",
		"-f", filepath.Join(dirpath, constant.FilenameQueriesSchemaSQL),
	)

	std, err := pgdump.CombinedOutput()
	if err != nil {
		return fmt.Errorf("write schema file: pg_dump: %q: %w", string(std), err)
	}

	return nil
}

// SchemaGo generates a schema.go file in the given directory using the
// database connection string and database schema name.
func SchemaGo(dirpath, dbconnection, dbschema string) error {
	// dirpathGenerationSpace represents the directory used to contain files during generation.
	//
	// The dirpathGenerationSpace directory is used when Jet generates files from the database.
	dirpathGenerationSpace := filepath.Join(
		dirpath, // (e.g., `.../queries/schema` from `.../queries/schema/schema.go`)
		constant.DirnameTempQueriesGenerationGo,
	)

	// Do not overwrite an existing directory from the user.
	if _, err := os.Stat(dirpathGenerationSpace); err == nil {
		return fmt.Errorf("warning: directory at must be deleted: %q", dirpathGenerationSpace)
	} else if errors.Is(err, os.ErrNotExist) {

	} else {
		return fmt.Errorf("error checking for directory space: %w", err)
	}

	// Generate the database schema models as Go types.
	generatorTemplate := genTemplate()
	if err := postgres.GenerateDSN(
		dbconnection,
		dbschema,
		dirpathGenerationSpace,
		generatorTemplate,
	); err != nil {
		return fmt.Errorf("jet: %w", err)
	}

	fmt.Println("Generated schema as models.")

	// Merge generated files to a single schema.go file.
	fileContentSchemas := [][]byte{
		[]byte("package " + constant.PkgNameSchemaGo + "\n\nimport \"github.com/switchupcb/jet/v2/postgres\""),
	}

	if err := filepath.WalkDir(dirpathGenerationSpace, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Do not merge generated files from the model directory.
			if filepath.Base(path) == "model" {
				return nil
			}

			// Do not attempt to merge directories with 0 files.
			fileCount, err := countDirFiles(path)
			if err != nil {
				return fmt.Errorf("directory file count: %w", err)
			}

			if fileCount == 0 {
				return nil
			}

			xstructOutput, err := xstruct(path, constant.PkgNameSchemaGo)
			if err != nil {
				return fmt.Errorf("xstruct called from %q: %w", path, err)
			}

			fileContentSchemas = append(fileContentSchemas, xstructOutput)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error flattening structs from generated SQL Go types: %w", err)
	}

	merger := newMerger(constant.PkgNameSchemaGo)
	for i := range fileContentSchemas {
		if err := merger.parseFile("", fileContentSchemas[i]); err != nil {
			return fmt.Errorf("merge: file_content_schema: %w\n\n%v", err, string(fileContentSchemas[i]))
		}
	}

	delete(merger.addedImports, "\"github.com/go-jet/jet/postgres\"")

	if err := merger.WriteToFile(filepath.Join(dirpath, constant.FilenameTemplateSchemaGo)); err != nil {
		return fmt.Errorf("merge: write: %w", err)
	}

	if err := os.RemoveAll(dirpathGenerationSpace); err != nil {
		return fmt.Errorf("merge: clean: %w", err)
	}

	return nil
}

// countDirFiles counts the number of non-directory files in a directory.
func countDirFiles(dirpath string) (int, error) {
	file, err := os.Open(dirpath)
	if err != nil {
		return 0, err //nolint:wrapcheck
	}

	defer file.Close()

	list, err := file.Readdirnames(-1)
	if err != nil {
		return 0, err //nolint:wrapcheck
	}

	var fileCount int

	for i := range list {
		fileInfo, err := os.Stat(filepath.Join(dirpath, list[i]))
		if err != nil {
			return 0, err //nolint:wrapcheck
		}

		if !fileInfo.IsDir() {
			fileCount++
		}
	}

	return fileCount, nil
}
