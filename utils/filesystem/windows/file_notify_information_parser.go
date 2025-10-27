//go:build windows
// +build windows

package windows

import (
	"unsafe"
	"golang.org/x/sys/windows"
	"path/filepath"
	"os"
	"time"
	"strings"
)

type FILE_NOTIFY_INFORMATION struct {
	NextEntryOffset uint32		// 4 bytes
	Action			uint32		// 4 bytes
	FileNameLength	uint32		// 4 bytes
	FileName		[1]uint16	// variable-length array, utf-16
}

type MetaData struct {
	ModTime 	time.Time
}

var cache map[string]MetaData = make(map[string]MetaData)
var debounceModTime time.Duration = 400 * time.Millisecond

func FileNotifyInformationParser(buffer [4096]byte, watchDirPath string, ignorePaths []string) ([]string, error) {
	var offset uint32 = 0
	actions := make([]string, 0, 1)
	oldNameActionIdx := -1
	for {
		entry := (*FILE_NOTIFY_INFORMATION)(unsafe.Pointer(&buffer[offset]))
		filename_b := (*[1<<10]uint16)(unsafe.Pointer(&entry.FileName[0]))[:entry.FileNameLength/2:entry.FileNameLength/2]
		filename := windows.UTF16ToString(filename_b)
		
		absPath, _ := filepath.Abs(filepath.Join(watchDirPath, filename))
		if !isValidEvent(absPath, ignorePaths) {
			return []string{}, nil
		}

		var action string;
		switch entry.Action {
		case windows.FILE_ACTION_ADDED:
			action = filename+" added"
		case windows.FILE_ACTION_REMOVED :
			action = filename+" removed"
		case windows.FILE_ACTION_MODIFIED:
			dupEvent, dupEventErr := isDuplicateModifiedEvent(absPath)
			if dupEventErr != nil {
				return []string{}, dupEventErr
			}
			if dupEvent {
				return []string{}, nil
			}
			action = filename+" modified"
		case windows.FILE_ACTION_RENAMED_OLD_NAME:
			action = filename+" renamed"
			oldNameActionIdx = len(actions)
		case windows.FILE_ACTION_RENAMED_NEW_NAME:
			if oldNameActionIdx > -1 {
				actions[oldNameActionIdx] += " to "+filename
				oldNameActionIdx = -1
			} else {
				action = "Renamed to "+filename
			}
		}
		if len(action) > 0 {
			actions = append(actions, action)
		}

		if entry.NextEntryOffset == 0 {
			break;
		}
		offset += entry.NextEntryOffset
	}
	return actions, nil
}

func isValidEvent(absPath string, ignorePaths []string) bool {
	for _, path := range ignorePaths {
		if strings.HasPrefix(absPath, path){
			return false
		}
	}
	return true
}

func isDuplicateModifiedEvent(absPath string) (bool, error) {
	stat, statErr := os.Stat(absPath)

	if statErr != nil {
		return false, statErr
	}

	currModTime := stat.ModTime()

	metadata, ok := cache[absPath]

	if !stat.IsDir() && (!ok || currModTime.Sub(metadata.ModTime) > debounceModTime) {
		cache[absPath] = MetaData{ ModTime: currModTime }
		return false, nil
	}

	return true, nil
}