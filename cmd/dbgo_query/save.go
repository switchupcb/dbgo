package query

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/switchupcb/dbgo/cmd/config"
)

// Save runs dbgo query save programmatically using the given filepath and YML.
func Save(abspath string, yml config.YML) (string, error) {
	filename := filepath.Base(abspath)
	sql_filename := filename + fileExtSQL
	sql_filepath := filepath.Join(yml.Generated.Input.Queries, sql_filename)

	fmt.Printf("SAVING QUERY %v from template at %v\n", sql_filename, abspath)

	// sqlcode is returned from an interpreted function which returns type-safe SQL in a string.
	sqlcode, err := interpretFunction(abspath)
	if err != nil {
		return "", fmt.Errorf("interpreter: %w", err)
	}

	// write output to an sql file with the same name as the interpreted file.
	if err := os.WriteFile(sql_filepath, []byte(sqlcode), writeFileMode); err != nil {
		return "", fmt.Errorf("error creating sql file: %w", err)
	}

	return fmt.Sprintf("%v QUERY SAVED from template at %v", sql_filename, abspath), nil
}
