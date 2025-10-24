//go:build windows
// +build windows

package filesystem

import (
	fs_types 	"github.com/blackistheneworange/filewatcher/utils/filesystem/types"
	"github.com/blackistheneworange/filewatcher/utils/filesystem/windows"
)

func GetOSFileSystemWatcher() fs_types.FileSystemWatcher {
	return &windows.WindowsFileSystemWatcher{}
}