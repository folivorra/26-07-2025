package downloader

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	urler "net/url"
	"os"
	"strings"
	"time"

	"github.com/folivorra/ziper/app"
	"github.com/google/uuid"
)

const destDir = "downloads"

type HTTPDownloader struct {
	client *http.Client
	a      *app.App
	logger *slog.Logger
}

var _ Downloader = (*HTTPDownloader)(nil)

func NewHTTPDownloader(a *app.App, logger *slog.Logger, timeout time.Duration) *HTTPDownloader {
	httpd := &HTTPDownloader{
		client: &http.Client{Timeout: timeout},
		a:      a,
		logger: logger,
	}

	httpd.a.RegisterCleanup(func(ctx context.Context) {
		if err := os.RemoveAll("downloads"); err != nil {
			httpd.logger.Warn("failed to remove downloads directory")
			return
		}
		httpd.logger.Info("removed downloads directory")
	})

	return httpd
}

func (d *HTTPDownloader) DownloadFile(url string, id uint64) error {
	parsedURL, err := urler.Parse(url)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}
	tokens := strings.Split(parsedURL.Path, "/")
	baseName := tokens[len(tokens)-1]
	if baseName == "" {
		baseName = "unnamed_download"
	}

	uuidStr := uuid.New().String()
	fileName := fmt.Sprintf("%s_%s", uuidStr, baseName)

	err = os.MkdirAll(fmt.Sprintf("%s/task-%d", destDir, id), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	filePath := fmt.Sprintf("%s/task-%d/%s", destDir, id, fileName)

	resp, err := d.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file, status: %s", resp.Status)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}
