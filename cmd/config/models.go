// Package config loads configuration data from an external file.
package config

// YML represents the first level of the YML file.
type YML struct {
	Abspath   string
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
	DB      DB     `yaml:"db"`
	Queries string `yaml:"queries"`
}

// DB represents generator input database properties of the YML file.
type DB struct {
	Connection string `yaml:"connection"`
	Schema     string `yaml:"schema"`
}

// Output represents generator output properties of the YML file.
type Output struct {
	Dpkg     string `yaml:"dpkg"`
	DBpkg    string `yaml:"dbpkg"`
	Template string `yaml:"template"`
}
