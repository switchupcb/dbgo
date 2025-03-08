package query

const (
	writeFileMode = 0644
	fileExtSQL    = ".sql"
	fileExtGo     = ".go"

	sqlGoDir                  = "go"
	queriesGoTemplatesDirname = "templates"

	newline    = '\n'
	colon      = ':'
	whitespace = ' '

	err_database_unspecified = "you must specify a database connection ('dbc') in the .yml configuration file"
)
