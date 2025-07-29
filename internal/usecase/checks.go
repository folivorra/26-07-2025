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

//func CanAddTask(activeTasks *atomic.Uint64, maxTasks uint64) bool {
//	for {
//		current := activeTasks.Load()
//		fmt.Println(current, maxTasks)
//		if current >= maxTasks {
//			return false
//		}
//		if activeTasks.CompareAndSwap(current, current+1) {
//			return true
//		}
//		// CAS-loop
//	}
//} // todo:
