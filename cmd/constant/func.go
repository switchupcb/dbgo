package constant

import (
	"os"
	"path/filepath"
)

// CopyFile copies a source file to a destination file.
func CopyFile(srcpath, dstpath string) error {
	src, err := os.ReadFile(srcpath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dstpath), FileModeWrite); err != nil {
		return err
	}

	if err := os.WriteFile(dstpath, src, FileModeWrite); err != nil {
		return err
	}

	return nil
}
