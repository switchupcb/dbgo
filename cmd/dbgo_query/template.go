package query

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-jet/jet/v2/generator/postgres"
	_ "github.com/lib/pq"
	"github.com/switchupcb/dbgo/cmd/config"
)

const (
	templateGoSchemaFilename = "schema.go"

	file_content_static = `

import . "github.com/go-jet/jet/v2/postgres"

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

	if yml.Generated.Input.DB.Connection == "" {
		return "", fmt.Errorf("you must specify a database connection ('dbc') in the .yml configuration file")
	}

	if yml.Generated.Input.DB.Schema == "" {
		yml.Generated.Input.DB.Schema = "public"
	}

	// Generate the database schema models as Go types.
	template_name := filepath.Base(abspath)
	sqlGoDirpath := filepath.Join(
		yml.Generated.Input.Queries, // queries
		queriesGoTemplatesDirname,   // templates
		template_name,               // template (name)
		sqlGoDir,                    // go
	)

	generatorTemplate := genTemplate()
	if err := postgres.GenerateDSN(
		yml.Generated.Input.DB.Connection,
		yml.Generated.Input.DB.Schema,
		sqlGoDirpath,
		generatorTemplate,
	); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	fmt.Println("Generated schema as models.")
	fmt.Println()

	// Merge generated files to a single schema.go file.
	file_content_schemas := [][]byte{
		[]byte("package " + template_name + "\n\n" + "import \"github.com/go-jet/jet/v2/postgres\""),
	}

	if err := filepath.WalkDir(sqlGoDirpath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if filepath.Base(path) == "model" {
				return nil
			}

			file_count, err := countDirFiles(path)
			if err != nil {
				return fmt.Errorf("directory file count: %w", err)
			}

			if file_count == 0 {
				return nil
			}

			// TODO: No guarantee of xstruct path yet
			xstruct := exec.Command("./xstruct", "-d", path, "-p", template_name, "-f", "-g")
			std, err := xstruct.CombinedOutput()
			if err != nil {
				return fmt.Errorf("xstruct called from %q: %v", path, string(std))
			}

			file_content_schemas = append(file_content_schemas, std)
		}

		return nil
	}); err != nil {
		return "", fmt.Errorf("error flattening structs from generated SQL Go types: %w", err)
	}

	merger := NewMerger(template_name)
	for i := range len(file_content_schemas) {
		if err := merger.parseFile("", file_content_schemas[i]); err != nil {
			return "", fmt.Errorf("merge: file_content_schema: %w\n\n%v", err, string(file_content_schemas[i]))
		}
	}

	delete(merger.addedImports, "\"github.com/go-jet/jet/postgres\"")

	templateSchemaFilepath := filepath.Join(
		yml.Generated.Input.Queries, // queries
		queriesGoTemplatesDirname,   // templates
		template_name,               // template (name)
		templateGoSchemaFilename,    // schema.go
	)

	if err := merger.WriteToFile(templateSchemaFilepath); err != nil {
		return "", fmt.Errorf("merge: write: %w", err)
	}

	if err := os.RemoveAll(sqlGoDirpath); err != nil {
		return "", fmt.Errorf("merge: clean: %w", err)
	}

	// Create the interpreted function file.
	file_content := []byte("package " + template_name + file_content_static)

	templateFilepath := filepath.Join(
		yml.Generated.Input.Queries, // queries
		queriesGoTemplatesDirname,   // templates
		template_name,               // template (name)
		template_name+fileExtGo,     // template.go
	)

	if err := os.WriteFile(templateFilepath, file_content, writeFileMode); err != nil {
		return "", fmt.Errorf("template: write: %w", err)
	}

	return fmt.Sprintf("added template %q to %v", template_name, filepath.Dir(templateFilepath)), nil
}

// countDirFiles counts the number of non-directory files in a directory.
func countDirFiles(dirpath string) (int, error) {
	file, err := os.Open(dirpath)
	if err != nil {
		return 0, err
	}

	defer file.Close()

	list, err := file.Readdirnames(-1)
	if err != nil {
		return 0, err
	}

	var file_count int
	for i := range len(list) {
		file_info, err := os.Stat(filepath.Join(dirpath, list[i]))
		if err != nil {
			return 0, err
		}

		if !file_info.IsDir() {
			file_count++
		}
	}

	return file_count, nil
}
