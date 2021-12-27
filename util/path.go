package util

import (
	"os"
	"path"

	"github.com/juju/errors"
)

// GenerateFilePath ...
func GenerateFilePath(dir, filename, ext string) (string, error) {
	if dir == "" || filename == "" {
		return "", errors.BadRequestf("util GenerateFilePath: Invalid parameter dir[%s] filename[%s]", dir, filename)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return "", errors.Annotate(err, "util GenerateFilePath")
		}
	}

	return path.Join(dir, filename+ext), nil
}
