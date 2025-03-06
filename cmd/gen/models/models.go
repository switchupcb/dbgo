package models

// Generator represents a code generator.
type Generator struct {
	Options GeneratorOptions // The custom options for the generator.
	Setpath string           // The filepath the setup file is located in.
	Outpath string           // The filepath the generated code is output to.
	Tempath string           // The filepath for the template used to generate code.
}

// GeneratorOptions represents options for a Generator.
type GeneratorOptions struct {
	Custom map[string]interface{} // The custom options of a generator.
}
