package dowloader

type Downloader interface {
	DownloadFile(url string, id uint64) (string, error)
}
