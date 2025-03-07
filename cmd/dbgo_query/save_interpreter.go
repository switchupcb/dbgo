package query

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/switchupcb/dbgo/cmd/dbgo_query/extract"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

const interpretedFunction = "sql.SQL"

// interpretFunction loads a template symbol from an interpreter.
func interpretFunction(dirpath string) (string, error) {
	files, err := os.ReadDir(dirpath)
	if err != nil {
		return "", fmt.Errorf("error loading template: %v\nIs the relative or absolute filepath set correctly?\n%w", dirpath, err)
	}

	// setup the interpreter
	_, err = os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("error loading default root directory for user-specific cached data. Is the GOCACHE set in `go env`?\n%w", err)
	}

	i := interp.New(interp.Options{GoPath: os.Getenv("GOPATH")})
	if err := i.Use(stdlib.Symbols); err != nil {
		return "", fmt.Errorf("error loading template stdlib libraries: %w", err)
	}

	// import jet symbols
	if err := i.Use(extract.Symbols); err != nil {
		return "", fmt.Errorf("error loading template imported symbols: %w", err)
	}

	// load the source
	for index := range files {
		srcpath := filepath.Join(dirpath, files[index].Name())
		src, err := os.ReadFile(srcpath)
		if err != nil {
			return "", fmt.Errorf("error evaluating template file src: %w", err)
		}

		if _, err := i.Eval(string(src)); err != nil {
			return "", fmt.Errorf("error evaluating template file: %w", err)
		}
	}

	// load the func from the interpreter
	v, err := i.Eval(interpretedFunction)
	if err != nil {
		return "", fmt.Errorf("error evaluating template function. Is it located in the file?\n%w", err)
	}

	fn, ok := v.Interface().(func() (string, error))
	if !ok {
		return "", errors.New("the template function `SQL` could not be type asserted. Is it a func() (string, error)?")
	}

	content, err := fn()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return content, nil
}
