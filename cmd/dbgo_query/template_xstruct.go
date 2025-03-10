package query

import (
	"fmt"

	xstructConfig "github.com/switchupcb/xstruct/cli/config"
	xstructGen "github.com/switchupcb/xstruct/cli/generator"
	xstructParser "github.com/switchupcb/xstruct/cli/parser"
	"golang.org/x/tools/imports"
)

// xstruct runs xstruct programmatically using the given path and package name.
func xstruct(dirpath, pkg string) ([]byte, error) {
	gen, err := xstructConfig.LoadFiles(dirpath)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err = xstructParser.Parse(gen, true, true); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	content := xstructGen.AstWriteDecls(pkg, gen.ASTDecls, gen.FuncDecls)

	// imports
	importsdata, err := imports.Process("", content, nil)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return importsdata, nil
}
