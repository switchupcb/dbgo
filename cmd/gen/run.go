package gen

import (
	"fmt"

	"github.com/switchupcb/dbgo/cmd/gen/config"
)

// Run runs dbgo programmatically using the given Environment's YMLPath.
func Run(ymlpath string) (string, error) {
	// The configuration file is loaded (.yml)
	_, err := config.LoadYML(ymlpath)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return "unimplemented", nil
}
