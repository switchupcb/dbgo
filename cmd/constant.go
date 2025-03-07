package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/switchupcb/dbgo/cmd/config"
)

const (
	flag_yml_name              = "yml"
	flag_yml_shorthand         = "y"
	flag_yml_usage             = `The path to the .yml flag used for code generation (from the current working directory)`
	flag_yml_usage_unspecified = "you must specify a .yml configuration file using --yml path/to/yml"

	queriesGoTemplatesDirname = "templates"
)

var (
	ymlFlag = new(string)
)

// parseFlagYML parses a "--yml" flag.
func parseYML() (*config.YML, error) {
	if ymlFlag == nil || *ymlFlag == "" {
		return nil, errors.New(flag_yml_usage_unspecified)
	}

	// The configuration file is loaded (.yml)
	yml, err := config.LoadYML(*ymlFlag)
	if err != nil {
		return nil, fmt.Errorf("error parsing yml: %w", err)
	}

	return yml, nil
}

// parseArgFilepath parses a filepath to return an absolute filepath to a template directory.
func parseArgFilepath(unknownpath string) (string, error) {
	if filepath.Ext(unknownpath) != "" {
		return "", fmt.Errorf("specified filepath is not a directory: %q", unknownpath)
	}

	// determine the actual filepath of the setup file.
	if !filepath.IsAbs(unknownpath) {
		if unknownpath, err := filepath.Abs(unknownpath); err != nil {
			return "", fmt.Errorf("error while determining the absolute file path of the specified template: %q\n%w", unknownpath, err)
		}
	}

	return unknownpath, nil
}
