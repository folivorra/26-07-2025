package archiver

type Archiver interface {
	ArchiveDirectory(dirPath string) error
}
