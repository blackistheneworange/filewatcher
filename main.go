package main

import (
	"flag"
	"github.com/blackistheneworange/filewatcher/utils"
	"github.com/blackistheneworange/filewatcher/utils/filesystem"
	"github.com/blackistheneworange/filewatcher/utils/logger"
	"github.com/blackistheneworange/filewatcher/utils/process"
)

var execCmd *string = flag.String("exec", "", "Command to execute")
var watchDirs *string = flag.String("watch", "", "Directory paths to watch")
var ignorePaths *string = flag.String("ignore", "", "Directories/files to ignore when watching")

func main() {
	flag.Parse()

	watcher, watcherErr := filesystem.GetFileSystemWatcher()
	if watcherErr != nil {
		logger.Fatal(logger.Format(0, watcherErr.Error()))
	}

	processManager, processManagerErr := process.GetProcessManager(execCmd)
	if processManagerErr != nil {
		logger.Fatal(logger.Format(0, processManagerErr.Error()))
	}
	process := processManager.CreateProcess(*execCmd)

	watchDirPaths, watchDirPathsErr := utils.GetWatchPaths(watchDirs)
	if watchDirPathsErr != nil {
		logger.Fatal(logger.Format(0, watchDirPathsErr.Error()))
	}

	ignorePaths, ignorePathsErr := utils.GetIgnorePaths(ignorePaths)
	if ignorePathsErr != nil {
		logger.Fatal(logger.Format(0, ignorePathsErr))
	}
	segregatedIgnorePaths := utils.SegregateIgnorePaths(ignorePaths, watchDirPaths)

	logger.Log("Listening for changes...")

	handleProcessError(process.StartProcess())

	ch := watcher.Subscribe()

	for idx, watchDirPath := range watchDirPaths {
		go watcher.Watch(watchDirPath, segregatedIgnorePaths[idx])
	}

	for {
		change := <-ch
		if change.Error != nil {
			logger.Fatal(logger.Format(' ', "watch error:", change.Error))
		}

		if len(change.Actions) > 0 {
			for _, action := range change.Actions {
				logger.Log(action)
			}
			logger.Log("Restarting...")

			handleProcessError(process.StartProcess())
		}
	}
}

func handleProcessError(err error) {
	if err != nil {
		logger.Log(err)
		logger.Log("Process crashed. Waiting for changes...")
	}
}
