package query

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/switchupcb/dbgo/cmd/config"
)

const (
	writeFileMode = 0644
	fileExtSQL    = ".sql"
	fileExtGo     = ".go"
)

// Add runs dbgo query add programmatically using the given filepath and YML.
func Add(abspath string, yml config.YML) (string, error) {
	// sqlcode is returned from an interpreted function which returns type-safe SQL in a string.
	sqlcode, err := interpretFunction(abspath)
	if err != nil {
		return "", fmt.Errorf("interpreter: %w", err)
	}

	// write output to an sql file with the same name as the interpreted file.
	filename := filepath.Base(abspath)
	sql_filename := filename[:len(filename)-len(fileExtGo)] + fileExtSQL
	sql_filepath := filepath.Join(yml.Generated.Input.Queries, sql_filename)

	if err := os.WriteFile(sql_filepath, []byte(sqlcode), writeFileMode); err != nil {
		return "", fmt.Errorf("error creating sql file: %w", err)
	}

	return fmt.Sprintf("%v query saved from %v", sql_filename, filename), nil
}
