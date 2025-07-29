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
	ID          uint64     `json:"id"`
	Status      TaskStatus `json:"status"`
	Files       []*File    `json:"files,omitempty"`
	ArchivePath string     `json:"archive_path,omitempty"`
	ArchiveURL  string     `json:"archive_url,omitempty"`
}

type File struct {
	Status FileStatus `json:"status"`
	URL    string     `json:"url"`
}
