package query

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/switchupcb/dbgo/cmd/constant"
	"github.com/switchupcb/dbgo/cmd/dbgo_query/extract"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

const interpretedFunctionName = "SQL"

// interpretFunction loads a template symbol from an interpreter.
func interpretFunction(dirpath string) (string, error) {
	_, err := os.ReadDir(dirpath)
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

	// load the source (in a specific order)
	if _, err := i.EvalPath(filepath.Join(dirpath, constant.FilenameTemplateSchemaGo)); err != nil {
		return "", fmt.Errorf("error compiling template schema.go file: %w", err)
	}

	if _, err := i.EvalPath(filepath.Join(dirpath, filepath.Base(dirpath)+constant.FileExtGo)); err != nil {
		return "", fmt.Errorf("error compiling template.go file: %w", err)
	}

	// load the func from the interpreter
	interpretedFunction := filepath.Base(dirpath) + "." + interpretedFunctionName

	v, err := i.Eval(interpretedFunction)
	if err != nil {
		return "", fmt.Errorf("error evaluating template function. Is it located in the file?\n%w", err)
	}

	fn, ok := v.Interface().(func() (string, error))
	if !ok {
		return "", errors.New("the template function `SQL` could not be type asserted. Is it a func() (string, error)?")
	}

	defer func() {
		if r := recover(); r != nil {
			r_msg, ok := r.(string)
			if !ok {
				fmt.Println("impossible recovery")
			}

			fmt.Printf("\t%v", r_msg)

			os.Exit(1)
		}
	}()

	content, err := fn()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return content, nil
}
