package model

type (
	TaskStatus string
	FileStatus string
)

const (
	TaskStatusAccepted   TaskStatus = "accepted"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"

	FileStatusAccepted    FileStatus = "accepted"
	FileStatusCompleted   FileStatus = "completed"
	FileStatusFailed      FileStatus = "failed"
	FileStatusInvalidType FileStatus = "invalid_type"
)

type Task struct {
	ID          uint64
	Status      TaskStatus
	Files       []*File
	ArchivePath string
	ArchiveURL  string
}

type File struct {
	Status FileStatus
	URL    string
}
