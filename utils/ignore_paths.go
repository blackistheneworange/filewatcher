package utils

import (
	"os"
	"path/filepath"
	"fmt"
	"errors"
	"strings"
)

func GetIgnorePaths(paths *string) ([]string, error) {
	if *paths == "" {
		return []string{}, nil
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
			return []string{}, errors.New(fmt.Sprintf("ignore path error: %s",pathErr.Error()))
		}
		_, statErr := os.Stat(path)
		if statErr != nil {
			return []string{}, errors.New(fmt.Sprintf("ignore path error: %s",statErr.Error()))
		}
		finalPathsList = append(finalPathsList, path)
	}

	return finalPathsList, nil
}

func SegregateIgnorePaths(ignorePaths []string, watchDirPaths []string) [][]string {
	segregatedList := make([][]string, len(watchDirPaths))

	for _, ignorePath := range ignorePaths {
		for idx, watchDirPath := range watchDirPaths {
			if strings.HasPrefix(ignorePath, watchDirPath) {
				segregatedList[idx] = append(segregatedList[idx], ignorePath)
				break;
			}
		}
	}

	return segregatedList
}