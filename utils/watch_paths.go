package utils

import (
	"os"
	"path/filepath"
	"fmt"
	"errors"
	"strings"
)

func GetWatchPaths(paths *string) ([]string, error) {
	if *paths == "" {
		return []string{}, errors.New("watch path error: No watch paths provided")
	}

	pathsList := strings.Split(*paths, ",")
	finalPathsList := make([]string, 0, len(pathsList))
	var pathErr error

	for _, path := range pathsList {
		path = strings.TrimSpace(path)
		if len(path) == 0 {
			continue
		}

		path, pathErr = filepath.Abs(path)
		if pathErr != nil {
			return []string{}, errors.New(fmt.Sprintf("watch path error: %s",pathErr.Error()))
		}
		stat, statErr := os.Stat(path)
		if statErr != nil {
			return []string{}, errors.New(fmt.Sprintf("watch path error: %s",statErr.Error()))
		}
		if !stat.IsDir() {
			return []string{}, errors.New(fmt.Sprintf("watch path error: Provided watch path %s does not point to a directory", path))
		}
		finalPathsList = append(finalPathsList, path)
	}

	return finalPathsList, nil
}