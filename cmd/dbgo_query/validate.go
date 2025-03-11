package query

import (
	"os"

	"github.com/switchupcb/dbgo/cmd/config"
	"github.com/switchupcb/dbgo/cmd/constant"
)

// validatedDatabaseConnection returns a validated Database Connection string from the given YML.
func validatedDatabaseConnection(yml config.YML) (string, error) {
	if yml.Generated.Input.DB.Connection == "" {
		return "", constant.ErrYMLDatabaseUnspecified
	}

	if yml.Generated.Input.DB.Connection[0] == constant.DatabaseConnectionEnvironmentVariableSymbol {
		yml.Generated.Input.DB.Connection = os.Getenv(yml.Generated.Input.DB.Connection[1:])

		if yml.Generated.Input.DB.Connection == "" {
			return "", constant.ErrYMLDatabaseUnspecified
		}
	}

	return yml.Generated.Input.DB.Connection, nil
}
