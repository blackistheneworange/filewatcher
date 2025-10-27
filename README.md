# FileWatcher
A simple utility to watch and restart on directory changes. Only supported in windows for now

[Latest release](https://github.com/blackistheneworange/filewatcher/releases/tag/v0.2.0)

## Usage

`filewatcher.exe -exec "node index.js" -watch "./path" -ignore "ignorethisfile.js, alsoignorethis.ts"`

### -exec
> The command to be executed which runs the intended application

### -watch
> The directory to be watched. This path should always point to a directory

### -ignore
> CSV of files and directories to ignore when watching


## Additional Information

> [!NOTE]
> All file/directory paths can be either relative or absolute.

> [!CAUTION]
> The executable file may be false flagged for virus by some antivirus softwares due to use of certain low level windows process apis.
