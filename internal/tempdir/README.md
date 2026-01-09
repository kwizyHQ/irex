# tempdir

This package provides a singleton temporary directory for the program lifecycle.

## Features
- Ensures only one temp directory exists (singleton pattern)
- Creates the temp directory on first use
- Cleans up the directory on program exit (SIGINT/SIGTERM) or manual deletion
- Methods to:
  - Copy files/folders from `embed.FS` or any `fs.FS` to the temp directory
  - Overwrite files in the temp directory
  - Delete files or the entire temp directory
  - Clean a subdirectory inside the temp directory (delete all contents)

## Usage

```go
import "internal/tempdir"

dir := tempdir.Get()
path := dir.Path()
err := dir.CopyFile(embedFS, "file.txt", "dest.txt")
err := dir.CopyFolder(embedFS, "folder", "dest-folder")
err := dir.OverwriteFile("file.txt", []byte("data"))
err := dir.DeleteFile("file.txt")
err := dir.Delete() // deletes the temp directory
err := dir.Clean("subfolder") // deletes all contents inside tempdir/subfolder
```

## Notes
- Only one temp directory is created and used throughout the program lifecycle.
- Directory is deleted on exit (SIGINT/SIGTERM) or by calling `Delete()`.
