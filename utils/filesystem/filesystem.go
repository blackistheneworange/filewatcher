package filesystem

import (
	"runtime"
	"errors"
	"fmt"
	fs_types	"github.com/blackistheneworange/filewatcher/utils/filesystem/types"
)

const OS_WINDOWS = "windows"

func GetFileSystemWatcher() (fs_types.FileSystemWatcher, error) {
	switch runtime.GOOS {
	case OS_WINDOWS :
		return GetOSFileSystemWatcher(), nil
	default :
		return nil, errors.New(fmt.Sprintf("%s os is not supported", runtime.GOOS))
	}
}

