package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/switchupcb/dbgo/cmd/config"
)

const (
	flagYMLName                  = "yml"
	flagYMLShorthand             = "y"
	flagYMLUsage                 = `The path to the .yml file used for code generation (from the current working directory).`
	flagYMLUsageErrorUnspecified = "you must specify a .yml configuration file using --yml path/to/yml"
)

var (
	flagYML = new(string)
)

// parseFlagYML parses a "--yml" flag.
func parseFlagYML() (*config.YML, error) {
	if flagYML == nil || *flagYML == "" {
		return nil, errors.New(flagYMLUsageErrorUnspecified)
	}

	// The configuration file is loaded (.yml)
	yml, err := config.LoadYML(*flagYML)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	return yml, nil
}

// parseArgFilepath parses a filepath argument to return an absolute filepath to a template directory.
func parseArgFilepath(unknownpath string) (string, error) {
	if filepath.Ext(unknownpath) != "" {
		return "", fmt.Errorf("specified filepath is not a directory: %q", unknownpath)
	}

	// determine the actual filepath of the setup file.
	if !filepath.IsAbs(unknownpath) {
		if unknownpath, err := filepath.Abs(unknownpath); err != nil {
			return "", fmt.Errorf("error determining the absolute file path of the specified template: %q\n%w", unknownpath, err)
		}
	}

	return unknownpath, nil
}
