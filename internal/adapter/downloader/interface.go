package downloader

type Downloader interface {
	DownloadFile(url string, id uint64) error
}
