package constant

import "errors"

// Constants defined by the program.
const (
	DatabaseConnectionEnvironmentVariableSymbol = '$'
	DatabaseSchemaNameDefault                   = "public"

	DirnameQueriesTemplates = "templates"
	DirnameQueriesSchema    = "schema"

	FilenameTemplateSchemaGo = "schema.go"
)

// Variables defined by the program.
var (
	// ErrYMLDatabaseUnspecified represents an error with the configuration file's database connection value.
	ErrYMLDatabaseUnspecified = errors.New("you must specify a database connection ('dbc') in the .yml configuration file")
)
