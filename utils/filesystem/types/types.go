package types

type WatchUpdate struct {
	Error		error
	Actions		[]string
}

type FileSystemWatcher interface {
	Watch(string, []string)
	Subscribe() <- chan WatchUpdate
}

type FileSystemWatchProps struct {
	WatchUpdateChannel	chan WatchUpdate
}