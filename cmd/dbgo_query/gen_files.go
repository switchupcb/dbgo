package query

const (
	file_name_dummy_sql    = "query.sql"
	file_content_dummy_sql = `-- name: dummy :one
SELECT current_timestamp;
	`

	file_name_sqlc_json    = "crud.json"
	file_content_sqlc_json = `{
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
  "sql": [
    {
      "schema": "schema.sql",
      "queries": "query.sql",
      "engine": "postgresql",
      "codegen": [
        {
          "out": "output",
          "plugin": "gen-crud",
          "options": {}
        }
      ]
    }
  ]
}`
)
