package model

type (
	TaskStatus string
	FileStatus string
)

const (
	TaskStatusAccepted   TaskStatus = "accepted"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"

	FileStatusAccepted         FileStatus = "accepted"
	FileStatusCompleted        FileStatus = "completed"
	FileStatusFailed           FileStatus = "failed"
	FileStatusInvalidURL       FileStatus = "invalid_url"
	FileStatusNotReachable     FileStatus = "not_reachable"
	FileStatusNotSupportedType FileStatus = "not_supported_type"
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
