package query

import (
	"errors"
	"fmt"
	"os"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

const interpretedFunction = "sql.SQL"

// interpretFunction loads a template symbol from an interpreter.
func interpretFunction(filepath string) (string, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("error loading template file: %v\nIs the relative or absolute filepath set correctly?\n%w", filepath, err)
	}

	// setup the interpreter
	_, err = os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("error loading template file. Is the GOCACHE set in `go env`?\n%w", err)
	}

	i := interp.New(interp.Options{GoPath: os.Getenv("GOPATH")})
	if err := i.Use(stdlib.Symbols); err != nil {
		return "", fmt.Errorf("error loading template stdlib libraries\n%w", err)
	}

	// load the source
	if _, err := i.Eval(string(file)); err != nil {
		return "", fmt.Errorf("error evaluating template file\n%w", err)
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
