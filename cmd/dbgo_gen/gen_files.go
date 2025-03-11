package gen

import (
	"fmt"
	"path/filepath"

	"github.com/switchupcb/dbgo/cmd/config"
	"github.com/switchupcb/dbgo/cmd/constant"
)

// fileContentSQLCYML returns the sqlc Generator YML file content.
func fileContentSQLCYML(yml config.YML) (string, error) {
	filepathQueriesSQL := filepath.Join(
		yml.Generated.Input.Queries,
		constant.DirnameTempQueriesGenerationSQLC,
		constant.FilenameQueriesCombinedSQL,
	)

	filepathSchemaSQL := filepath.Join(
		yml.Generated.Input.Queries,
		constant.DirnameQueriesSchema,
		constant.FilenameQueriesSchemaSQL,
	)

	pkg := filepath.Base(yml.Generated.Output.DBpkg)

	// find the relative path of each file (relative to sqlc.yaml)
	filepathConfig := filepath.Join(
		yml.Generated.Input.Queries,
		constant.DirnameTempQueriesGenerationSQLC,
	)

	relativeQueriesPath, err := filepath.Rel(filepathConfig, filepathQueriesSQL)
	if err != nil {
		return "", fmt.Errorf("queries: %w", err)
	}

	relativeSchemaPath, err := filepath.Rel(filepathConfig, filepathSchemaSQL)
	if err != nil {
		return "", fmt.Errorf("queries: %w", err)
	}

	relativeOutputPath, err := filepath.Rel(filepathConfig, yml.Generated.Output.DBpkg)
	if err != nil {
		return "", fmt.Errorf("queries: %w", err)
	}

	return fmt.Sprintf(`version: "2"
sql:
  - engine: "postgresql"
    queries: "%s"
    schema: "%s"
    gen:
      go:
        package: "%s"
        out: "%s"
        sql_package: "pgx/v5"`,
		filepath.ToSlash(relativeQueriesPath), // queries
		filepath.ToSlash(relativeSchemaPath),  // schema
		pkg,                                   // package
		filepath.ToSlash(relativeOutputPath),  // out
	), nil
}
