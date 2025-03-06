// Package config loads configuration data from an external file.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/switchupcb/dbgo/cmd/gen/models"
	"gopkg.in/yaml.v3"
)

// LoadYML loads a .yml configuration file into a Generator.
func LoadYML(relativepath string) (*models.Generator, error) {
	file, err := os.ReadFile(relativepath)
	if err != nil {
		return nil, fmt.Errorf("the specified .yml filepath doesn't exist: %v\n%w", relativepath, err)
	}

	var yml YML
	if err := yaml.Unmarshal(file, &yml); err != nil {
		return nil, fmt.Errorf("an error occurred unmarshalling the .yml file\n%w", err)
	}

	gen := ParseYML(yml)

	// determine the actual filepath of the loader.
	absloadpath, err := filepath.Abs(relativepath)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while determining the absolute file path of the loader file\n%v", relativepath)
	}

	// determine the actual filepath of the setup.go file.
	gen.Setpath = filepath.Join(filepath.Dir(absloadpath), gen.Setpath)

	// determine the actual filepath of the template file (if provided).
	if gen.Tempath != "" {
		gen.Tempath = filepath.Join(filepath.Dir(absloadpath), gen.Tempath)
	}

	// determine the actual filepath of the output file.
	gen.Outpath = filepath.Join(filepath.Dir(absloadpath), gen.Outpath)

	return gen, nil
}

// ParseYML parses a YML into a Generator.
func ParseYML(yml YML) *models.Generator {
	return &models.Generator{
		Options: models.GeneratorOptions{
			Custom: yml.Options,
		},
		Setpath: "",
		Outpath: "",
		Tempath: yml.Generated.Output.Template,
	}
}
