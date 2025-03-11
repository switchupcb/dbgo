package query

import (
	"os"
	"path/filepath"

	"github.com/switchupcb/dbgo/cmd/constant"
)

// copyFile copies a source file to a destination file.
func copyFile(srcpath, dstpath string) error {
	src, err := os.ReadFile(srcpath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dstpath), constant.FileModeWrite); err != nil {
		return err
	}

	if err := os.WriteFile(dstpath, src, constant.FileModeWrite); err != nil {
		return err
	}

	return nil
}
