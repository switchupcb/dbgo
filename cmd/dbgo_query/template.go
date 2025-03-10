package query

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/switchupcb/dbgo/cmd/config"
	"github.com/switchupcb/dbgo/cmd/constant"
	"github.com/switchupcb/jet/v2/generator/postgres"
)

const (
	dirnameJetGenerated = "go"

	interpretedFileContentStatic = `

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

// Template runs dbgo query template programmatically using the given filepath and YML.
func Template(abspath string, yml config.YML) (string, error) {
	var err error

	yml.Generated.Input.DB.Connection, err = validatedDatabaseConnection(yml)
	if err != nil {
		return "", err
	}

	if yml.Generated.Input.DB.Schema == "" {
		yml.Generated.Input.DB.Schema = constant.DatabaseSchemaNameDefault
	}

	templateName := filepath.Base(abspath)

	fmt.Printf("ADDING TEMPLATE %v to %v\n\n", templateName, abspath)

	// Generate the database schema models as Go types.
	sqlGoDirpath := filepath.Join(
		yml.Generated.Input.Queries,      // queries
		constant.DirnameQueriesTemplates, // templates
		templateName,                     // template (name)
		dirnameJetGenerated,              // go
	)

	generatorTemplate := genTemplate()
	if err := postgres.GenerateDSN(
		yml.Generated.Input.DB.Connection,
		yml.Generated.Input.DB.Schema,
		sqlGoDirpath,
		generatorTemplate,
	); err != nil {
		return "", fmt.Errorf("jet: %w", err)
	}

	fmt.Println("Generated schema as models.")
	fmt.Println()

	// Merge generated files to a single schema.go file.
	fileContentSchemas := [][]byte{
		[]byte("package " + templateName + "\n\n" + "import \"github.com/switchupcb/jet/v2/postgres\""),
	}

	if err := filepath.WalkDir(sqlGoDirpath, func(path string, d fs.DirEntry, err error) error {
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

			xstructOutput, err := xstruct(path, templateName)
			if err != nil {
				return fmt.Errorf("xstruct called from %q: %w", path, err)
			}

			fileContentSchemas = append(fileContentSchemas, xstructOutput)
		}

		return nil
	}); err != nil {
		return "", fmt.Errorf("error flattening structs from generated SQL Go types: %w", err)
	}

	merger := newMerger(templateName)
	for i := range fileContentSchemas {
		if err := merger.parseFile("", fileContentSchemas[i]); err != nil {
			return "", fmt.Errorf("merge: file_content_schema: %w\n\n%v", err, string(fileContentSchemas[i]))
		}
	}

	delete(merger.addedImports, "\"github.com/go-jet/jet/postgres\"")

	templateSchemaFilepath := filepath.Join(
		yml.Generated.Input.Queries,       // queries
		constant.DirnameQueriesTemplates,  // templates
		templateName,                      // template (name)
		constant.FilenameTemplateSchemaGo, // schema.go
	)

	if err := merger.WriteToFile(templateSchemaFilepath); err != nil {
		return "", fmt.Errorf("merge: write: %w", err)
	}

	if err := os.RemoveAll(sqlGoDirpath); err != nil {
		return "", fmt.Errorf("merge: clean: %w", err)
	}

	// Create the interpreted function file.
	fileContent := []byte("package " + templateName + interpretedFileContentStatic)

	templateFilepath := filepath.Join(
		yml.Generated.Input.Queries,      // queries
		constant.DirnameQueriesTemplates, // templates
		templateName,                     // template (name)
		templateName+constant.FileExtGo,  // template.go
	)

	if _, err := os.Stat(templateFilepath); err == nil {

	} else if errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(templateFilepath, fileContent, constant.FileModeWrite); err != nil {
			return "", fmt.Errorf("template: write: %w", err)
		}
	} else {
		return "", fmt.Errorf("error checking for template file space: %w", err)
	}

	return fmt.Sprintf("ADDED TEMPLATE %q to %v", templateName, filepath.Dir(templateFilepath)), nil
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

	var file_count int

	for i := range list {
		file_info, err := os.Stat(filepath.Join(dirpath, list[i]))
		if err != nil {
			return 0, err //nolint:wrapcheck
		}

		if !file_info.IsDir() {
			file_count++
		}
	}

	return file_count, nil
}
