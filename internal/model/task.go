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

	Max
)

type Task struct {
	ID          uint64     `json:"id"`
	Status      TaskStatus `json:"status"`
	Files       []*File    `json:"urls"`
	ArchivePath string     `json:"archive_path,omitempty"`
}

type File struct {
	Status FileStatus `json:"status"`
	URL    string     `json:"url"`
}
