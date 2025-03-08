// Package config loads configuration data from an external file.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadYML loads a .yml configuration file into a Generator.
func LoadYML(relativepath string) (*YML, error) {
	file, err := os.ReadFile(relativepath)
	if err != nil {
		return nil, fmt.Errorf("the specified .yml filepath doesn't exist: %v\n%w", relativepath, err)
	}

	yml := new(YML)
	if err := yaml.Unmarshal(file, yml); err != nil {
		return nil, fmt.Errorf("error occurred unmarshalling the .yml file\n%w", err)
	}

	// determine the actual filepath of the setup file.
	absloadpath, err := filepath.Abs(relativepath)
	if err != nil {
		return nil, fmt.Errorf("error occurred while determining the absolute file path of the setup file\n%v", relativepath)
	}

	yml.abspath = absloadpath

	// determine the actual filepath of the domain package and queries directory.
	yml.Generated.Input.Dpkg = filepath.Join(filepath.Dir(absloadpath), yml.Generated.Input.Dpkg)
	yml.Generated.Input.Queries = filepath.Join(filepath.Dir(absloadpath), yml.Generated.Input.Queries)

	// determine the actual filepath of the domain package,  queries directory, and template.
	yml.Generated.Output.Dpkg = filepath.Join(filepath.Dir(absloadpath), yml.Generated.Output.Dpkg)
	yml.Generated.Output.DBpkg = filepath.Join(filepath.Dir(absloadpath), yml.Generated.Output.DBpkg)

	// determine the actual filepath of the template file (if provided).
	if yml.Generated.Output.Template != "" {
		yml.Generated.Output.Template = filepath.Join(filepath.Dir(absloadpath), yml.Generated.Output.Template)
	}

	return yml, nil
}
