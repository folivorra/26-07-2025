package dowloader

import (
	"fmt"
	"io"
	"net/http"
	urler "net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const destDir = "downloads"

type HTTPDownloader struct {
	client *http.Client
}

func NewHTTPDownloader(timeout time.Duration) *HTTPDownloader {
	return &HTTPDownloader{
		client: &http.Client{Timeout: timeout},
	}
}

func (d *HTTPDownloader) DownloadFile(url string, id uint64) (string, error) {
	// Парсим URL и извлекаем имя файла
	parsedURL, err := urler.Parse(url)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}
	tokens := strings.Split(parsedURL.Path, "/")
	fileName := tokens[len(tokens)-1]
	if fileName == "" {
		fileName = "unnamed_download"
	}

	// Создаём папку если нужно
	err = os.MkdirAll(fmt.Sprintf("%s/task-%d", destDir, id), os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create destination directory: %w", err)
	}

	filePath := filepath.Join(destDir, fileName)

	// HTTP GET
	resp, err := d.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file, status: %s", resp.Status)
	}

	// Создаём файл
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Пишем содержимое
	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	fmt.Printf("Downloaded file %s with length %d\n", fileName, written)

	return filePath, nil
}
