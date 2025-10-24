package utils

import (
	"os"
	"path/filepath"
	"fmt"
	"errors"
)

func GetWatchPath(path *string) (string, error) {
	if *path == "" {
		return "", errors.New("watch path error: No watch path provided")
	}

	absPath, absPathErr := filepath.Abs(*path)
	if absPathErr != nil {
		return "", absPathErr
	}

	watchDirInfo, statErr := os.Stat(absPath)
	if statErr != nil {
		return "", errors.New(fmt.Sprintf("watch path error: %s",statErr.Error()))
	}
	if !watchDirInfo.IsDir() {
		return "", errors.New(fmt.Sprintf("watch path error: Provided watch path %s does not point to a directory", *path))
	}

	return absPath, nil
}