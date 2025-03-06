// Package config loads configuration data from an external file.
package config

// YML represents the first level of the YML file.
type YML struct {
	Generated Generated              `yaml:"generated"`
	Options   map[string]interface{} `yaml:"custom"`
}

// Generated represents generated properties of the YML file.
type Generated struct {
	Input  Input  `yaml:"input"`
	Output Output `yaml:"output"`
}

// Input represents generator input properties of the YML file.
type Input struct {
	Dpkg    string `yaml:"dpkg"`
	Dbc     string `yaml:"dbc"`
	Queries string `yaml:"queries"`
}

// Output represents generator output properties of the YML file.
type Output struct {
	Dpkg     string `yaml:"dpkg"`
	DBpkg    string `yaml:"dbpkg"`
	Template string `yaml:"template"`
}
