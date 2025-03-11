package query

import (
	"fmt"
	"path/filepath"
)

const (
	dirnameGeneratedQueriesOutput = "output"
	filenameGeneratedQueriesSQL   = "query.sql"

	filenamePlaceholderSQL    = "placeholder.sql"
	fileContentPlaceholderSQL = `-- name: placeholder :one
SELECT current_timestamp;
	`

	filenameSQLCJSON = "crud.json"
)

// fileContentSQLCJSON returns the sqlc CRUD Generator JSON file content.
func fileContentSQLCJSON(schemapath string) string {
	return `{
  "version": "2",
  "plugins": [
    {
      "name": "gen-crud",
      "wasm": {
        "url": "https://github.com/kaashyapan/sqlc-gen-crud/releases/download/v1.0.1/sqlc-gen-crud_1.0.1.wasm",
        "sha256": "1a8146b30585882a8104d2ddcbfef0438b953fff08e74e7b90a9bf3d7bb2764c"
      }
    }
  ],
  "sql": [` +
		fmt.Sprintf(`
    {
      "schema": "%s",
      "queries": "%s",
      "engine": "postgresql",
      "codegen": [
        {
          "out": "%s",
          "plugin": "gen-crud",
          "options": {}
        }
      ]
    }
    `, filepath.ToSlash(schemapath), filenamePlaceholderSQL, dirnameGeneratedQueriesOutput) +
		`
  ]
}`
}
