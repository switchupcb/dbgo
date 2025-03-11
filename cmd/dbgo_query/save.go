package query

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/switchupcb/dbgo/cmd/config"
	"github.com/switchupcb/dbgo/cmd/constant"
)

// Save runs dbgo query save programmatically using the given template name and YML.
func Save(name string, yml config.YML) error {
	templateFilepath := filepath.Join(
		yml.Generated.Input.Queries,      // queries
		constant.DirnameQueriesTemplates, // templates
		name,                             // template (name)
	)

	if _, err := os.Stat(templateFilepath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf( //nolint:stylecheck // ST1005
				`template was not found at %q
	You can use 'db query template %v' to create a template.
	Then, run 'db query save %v' again to save the SQL output to an SQL file.`,
				templateFilepath, name, name,
			)
		} else {
			return fmt.Errorf("error checking template file space: %w", err)
		}
	}

	// sqlcode is returned from an interpreted function which returns type-safe SQL in a string.
	sqlcode, err := interpretFunction(templateFilepath)
	if err != nil {
		return fmt.Errorf("interpreter: %w", err)
	}

	if len(sqlcode) > 1 && sqlcode[0] == constant.Newline {
		sqlcode = sqlcode[1:]
	}

	// write output to an sql file with the same name as the interpreted file.
	filepathSQL := filepath.Join(yml.Generated.Input.Queries, name+constant.FileExtSQL)
	if err := os.WriteFile(filepathSQL, []byte(sqlcode), constant.FileModeWrite); err != nil {
		return fmt.Errorf("error creating sql file: %w", err)
	}

	return nil
}
