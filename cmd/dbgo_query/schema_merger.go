package query

// MIT License
//
// Copyright (c) 2023 Dmitrii Barsukov
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"sort"
)

type merger struct {
	tree         *ast.File
	addedImports map[string]ast.Spec
	addedConsts  map[string]ast.Spec
	addedTypes   map[string]ast.Spec
	addedVars    map[string]ast.Spec
	addedFunc    map[string]*ast.FuncDecl
	specialFunc  map[string]bool
}

func newMerger(pkg string) *merger {
	merger := &merger{
		tree: &ast.File{
			Name: ast.NewIdent(pkg),
		},
		addedImports: make(map[string]ast.Spec),
		addedConsts:  make(map[string]ast.Spec),
		addedTypes:   make(map[string]ast.Spec),
		addedVars:    make(map[string]ast.Spec),
		addedFunc:    make(map[string]*ast.FuncDecl),
		specialFunc:  map[string]bool{`init`: true, `main`: true},
	}

	return merger
}

func (m *merger) parseFile(path string, src any) error {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, path, src, 0)
	if err != nil {
		return err
	}

	for i, decl := range file.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok {
			if gen.Tok == token.PACKAGE {
				file.Decls = append(file.Decls[:i], file.Decls[i+1:]...)
				break
			}
		}
	}

	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			m.parseGenDecl(decl)
		case *ast.FuncDecl:
			name := decl.Name.Name
			if _, ok := m.specialFunc[name]; !ok {
				name = fset.Position(decl.Pos()).String()
			}
			m.addedFunc[name] = decl
		}
	}

	return nil
}

func (m *merger) parseGenDecl(decl *ast.GenDecl) {
	switch decl.Tok {
	case token.IMPORT:
		for _, spec := range decl.Specs {
			if v, ok := spec.(*ast.ImportSpec); ok {
				m.addedImports[v.Path.Value] = spec
			}
		}
	case token.CONST:
		for _, spec := range decl.Specs {
			if v, ok := spec.(*ast.ValueSpec); ok {
				for _, name := range v.Names {
					m.addedConsts[name.Name] = spec
				}
			}
		}
	case token.TYPE:
		for _, spec := range decl.Specs {
			if t, ok := spec.(*ast.TypeSpec); ok {
				m.addedTypes[t.Name.Name] = spec
			}
		}
	case token.VAR:
		for _, spec := range decl.Specs {
			if v, ok := spec.(*ast.ValueSpec); ok {
				for _, name := range v.Names {
					m.addedVars[name.Name] = spec
				}
			}
		}
	}
}

func (m *merger) buildGenDecl() {
	var specs []ast.Spec

	specs = make([]ast.Spec, 0, len(m.addedImports))
	for _, spec := range m.addedImports {
		specs = append(specs, spec)
	}
	if len(specs) > 0 {
		m.tree.Decls = append(m.tree.Decls, &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: specs,
		})
	}

	specs = make([]ast.Spec, 0, len(m.addedConsts))
	for _, spec := range m.addedConsts {
		specs = append(specs, spec)
	}
	if len(specs) > 0 {
		m.tree.Decls = append(m.tree.Decls, &ast.GenDecl{
			Tok:   token.CONST,
			Specs: specs,
		})
	}

	specs = make([]ast.Spec, 0, len(m.addedVars))
	for _, spec := range m.addedVars {
		specs = append(specs, spec)
	}
	if len(specs) > 0 {
		m.tree.Decls = append(m.tree.Decls, &ast.GenDecl{
			Tok:   token.VAR,
			Specs: specs,
		})
	}

	specs = make([]ast.Spec, 0, len(m.addedTypes))
	for _, spec := range m.addedTypes {
		specs = append(specs, spec)
	}
	if len(specs) > 0 {
		m.tree.Decls = append(m.tree.Decls, &ast.GenDecl{
			Tok:   token.TYPE,
			Specs: specs,
		})
	}
}

func (m *merger) sortAddedFuncs() []*ast.FuncDecl {
	keys := make([]string, 0, len(m.addedFunc))
	for k := range m.addedFunc {
		if _, ok := m.specialFunc[k]; ok {
			continue
		}
		keys = append(keys, k)
	}

	sort.Strings(keys)

	sortedFuncs := make([]*ast.FuncDecl, len(keys))
	for i, k := range keys {
		sortedFuncs[i] = m.addedFunc[k]
	}

	for k := range m.specialFunc {
		if _, ok := m.addedFunc[k]; !ok {
			continue
		}
		sortedFuncs = append(sortedFuncs, m.addedFunc[k])
	}

	return sortedFuncs
}

func (m *merger) WriteToFile(sourceName string) error {
	source, err := os.Create(sourceName)
	if err != nil {
		return err
	}
	defer func(source *os.File) {
		_ = source.Close()
	}(source)

	m.buildGenDecl()

	for _, decl := range m.sortAddedFuncs() {
		m.tree.Decls = append(m.tree.Decls, decl)
	}

	if err := printer.Fprint(source, token.NewFileSet(), m.tree); err != nil {
		return err
	}

	return nil
}
