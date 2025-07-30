package usecase

import (
	"path"
)

func IsAllowedFileType(url string) bool {
	return path.Ext(url) == ".pdf" || path.Ext(url) == ".jpeg"
}

func CanAddFileInTask(activeFiles uint64, maxFiles uint64) bool {
	return activeFiles < maxFiles
}
