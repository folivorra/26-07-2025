package archiver

type Archiver interface {
	ArchiveDirectory(dirPath, zipPath string) error
}
