package main

import (
	"flag"
	"github.com/blackistheneworange/filewatcher/utils"
	"github.com/blackistheneworange/filewatcher/utils/filesystem"
	"github.com/blackistheneworange/filewatcher/utils/logger"
	"github.com/blackistheneworange/filewatcher/utils/process"
)

var execCmd *string = flag.String("exec", "", "Command to execute")
var watchDir *string = flag.String("watch", "", "Directory path to watch")
var ignorePaths *string = flag.String("ignore", "", "Directories/files to ignore when watching")

func main() {
	flag.Parse()

	watcher, watcherErr := filesystem.GetFileSystemWatcher()
	if watcherErr != nil {
		panic(logger.Format(0, watcherErr.Error()))
	}

	processManager, processManagerErr := process.GetProcessManager(execCmd)
	if processManagerErr != nil {
		panic(logger.Format(0, processManagerErr.Error()))
	}
	process := processManager.CreateProcess(*execCmd)

	watchDirPath, watchDirPathErr := utils.GetWatchPath(watchDir)
	if watchDirPathErr != nil {
		panic(logger.Format(0, watchDirPathErr.Error()))
	}

	ignorePaths, ignorePathsErr := utils.GetIgnorePaths(ignorePaths, watchDirPath)
	if ignorePathsErr != nil {
		panic(logger.Format(0, ignorePathsErr))
	}

	logger.Log("Listening for changes...")

	handleProcessError(process.StartProcess())

	ch := watcher.Subscribe()
	go watcher.Watch(watchDirPath, ignorePaths)

	for {
		change := <-ch
		if change.Error != nil {
			panic(logger.Format(' ', "watch error:", change.Error))
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
		logger.Log("Process crashed. Waiting for changes...")
	}
}
