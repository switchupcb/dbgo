package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/switchupcb/dbgo/cmd/config"
)

const (
	flag_yml_name              = "yml"
	flag_yml_usage             = `The path to the .yml flag used for code generation (from the current working directory)`
	flag_yml_usage_unspecified = "you must specify a .yml configuration file using -yml path/to/yml"

	flag_filepath_usage_unspecified = "you must specify a .go filename"
)

// parseFlagYML parses a "-yml" flag.
func parseYML(ymlFlag *pflag.Flag) (*config.YML, error) {
	// todo: https://github.com/spf13/cobra/issues/2250
	if ymlFlag == nil || ymlFlag.Value.String() == "" {
		return nil, errors.New(flag_yml_usage_unspecified)
	}

	// The configuration file is loaded (.yml)
	yml, err := config.LoadYML(ymlFlag.Value.String())
	if err != nil {
		return nil, fmt.Errorf("error parsing yml: %w", err)
	}

	return yml, nil
}

const fileExtGo = ".go"

// parseArgFilepath parses a filepath for the absolute path to .go file.
func parseArgFilepath(unknownpath string) (string, error) {
	if filepath.Ext(unknownpath) != fileExtGo {
		return "", fmt.Errorf("specified filepath is not a .go file: %q", unknownpath)
	}

	// determine the actual filepath of the setup file.
	if !filepath.IsAbs(unknownpath) {
		if unknownpath, err := filepath.Abs(unknownpath); err != nil {
			return "", fmt.Errorf("error while determining the absolute file path of the specified .go file: %q\n%w", unknownpath, err)
		}
	}

	return unknownpath, nil
}
