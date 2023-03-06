package shared

import (
	"io/fs"
	"os"
)

func ReadDir(dirPath string) []fs.DirEntry {
	files, _ := os.ReadDir(dirPath)
	return files
}
