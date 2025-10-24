//go:build !windows
// +build !windows

package filesystem

import (
	fs_types 	"github.com/blackistheneworange/filewatcher/utils/filesystem/types"
)

func GetOSFileSystemWatcher() fs_types.FileSystemWatcher {
	return nil
}