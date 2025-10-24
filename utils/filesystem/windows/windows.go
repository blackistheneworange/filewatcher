//go:build windows
// +build windows

package windows

import (
	"golang.org/x/sys/windows"
	"github.com/blackistheneworange/filewatcher/utils/logger"
	fs_types "github.com/blackistheneworange/filewatcher/utils/filesystem/types"
)

type WatchUpdate = fs_types.WatchUpdate

func getHandle(path string) (windows.Handle, error) {
	var finalPath *uint16 = windows.StringToUTF16Ptr(path)

	// https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-createfilea
	return windows.CreateFile(
		finalPath,
		windows.FILE_LIST_DIRECTORY,	// desired access - purpose is to get the contents within dir									
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,		// mode - access provided to other processes when the file is open in this program
		nil,
		windows.OPEN_EXISTING,		// create mode - expect the file path to exist always
		windows.FILE_FLAG_BACKUP_SEMANTICS,		// attributes - FILE_FLAG_BACKUP_SEMANTICS required to get handle for a directory
		0,
	)
}

func closeHandle(handle windows.Handle) {
	windows.CloseHandle(handle)
}

func (watcher *WindowsFileSystemWatcher) Watch(watchDirPath string, ignorePaths []string) {
	handle, handleErr := getHandle(watchDirPath)
	if handleErr != nil {
		panic(logger.Format(' ', "filesystem handle error:",handleErr))
	}
	defer closeHandle(handle)

	for {
		// https://learn.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-readdirectorychangesw
		err := windows.ReadDirectoryChanges(
			handle,
			&watcher.Buffer[0],
			uint32(len(watcher.Buffer)),
			true,
			windows.FILE_NOTIFY_CHANGE_LAST_WRITE|windows.FILE_NOTIFY_CHANGE_FILE_NAME|windows.FILE_NOTIFY_CHANGE_DIR_NAME,
			&watcher.BytesReturned,
			nil,
			0,
		)

		if watcher.WatchUpdateChannel != nil {
			if err != nil {
				watcher.WatchUpdateChannel <- WatchUpdate{ Error: err }
			} else {
				actions, parseErr := FileNotifyInformationParser(watcher.Buffer, watchDirPath, ignorePaths)
				
				if(parseErr != nil) {
					watcher.WatchUpdateChannel <- WatchUpdate{ Error: parseErr }
				} else {
					watcher.WatchUpdateChannel <- WatchUpdate{ Error: nil, Actions: actions }
				}
			}
		}
	}
}

func (watcher *WindowsFileSystemWatcher) Subscribe() <- chan WatchUpdate {
	ch := make(chan WatchUpdate)
	watcher.WatchUpdateChannel = ch

	return ch
}