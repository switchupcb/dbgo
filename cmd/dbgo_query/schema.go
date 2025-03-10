package query

import (
	"path/filepath"

	"github.com/switchupcb/dbgo/cmd/config"
	"github.com/switchupcb/dbgo/cmd/constant"
)

// Schema runs dbgo query schema programmatically using the given YML.
func Schema(yml config.YML) (string, error) {
	queriesSchemaDir := filepath.Join(yml.Generated.Input.Queries, constant.DirnameQueriesSchema)

	return "Generated schema file(s)." + queriesSchemaDir, nil
}
