package utils

import (
	"os"
	"path/filepath"
	"fmt"
	"errors"
	"strings"
)

func GetIgnorePaths(paths *string, parentPath string) ([]string, error) {
	if *paths == "" {
		return []string{}, nil
	}

	pathsList := strings.Split(*paths, ",")
	var pathErr error

	for idx, path := range pathsList {
		path = strings.TrimSpace(path)
		path, pathErr = filepath.Abs(filepath.Join(parentPath, path))
		if pathErr != nil {
			return []string{}, errors.New(fmt.Sprintf("ignore path error: %s",pathErr.Error()))
		}
		_, statErr := os.Stat(path)
		if statErr != nil {
			return []string{}, errors.New(fmt.Sprintf("ignore path error: %s",statErr.Error()))
		}
		pathsList[idx] = path
	}

	return pathsList, nil
}