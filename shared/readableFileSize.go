package shared

import "fmt"

func ReadablefileSize(size int64) string {
	var fileSize float64 = float64(size)
	var formatedSize string
	if fileSize < 122400 {
		// return kb
		result := fileSize / 1024
		formatedSize = fmt.Sprintf("%.2f KB", result)
	} else if fileSize < 1080000000 {
		// return mb
		result := fileSize / (1024 * 1024)
		formatedSize = fmt.Sprintf("%.2f MB", result)
	} else {
		// return gb
		result := fileSize / (1024 * 1024 * 1024)
		formatedSize = fmt.Sprintf("%.2f GB", result)
	}
	return formatedSize
}
