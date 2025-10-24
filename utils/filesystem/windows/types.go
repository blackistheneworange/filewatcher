//go:build windows
// +build windows

package windows

import (
	fs_types "github.com/blackistheneworange/filewatcher/utils/filesystem/types"
)

type WindowsFileSystemWatcher struct {
	fs_types.FileSystemWatchProps
	Buffer 				[4096]byte
	BytesReturned		uint32
}