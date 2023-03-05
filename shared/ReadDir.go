package shared

import (
	"io/fs"
	"os"
)
func ReadDir(dirPath string)([]fs.DirEntry){
	files, err := os.ReadDir(dirPath)
	CheckError(err)
	return files 
}