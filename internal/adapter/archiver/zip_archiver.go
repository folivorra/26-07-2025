package archiver

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/folivorra/ziper/app"
)

type ZipArchiver struct {
	a      *app.App
	logger *slog.Logger
}

func NewZipArchiver(a *app.App, logger *slog.Logger) *ZipArchiver {
	za := &ZipArchiver{
		a:      a,
		logger: logger,
	}

	za.a.RegisterCleanup(func(ctx context.Context) {
		if err := os.RemoveAll("archives"); err != nil {
			za.logger.Warn("failed to remove old archives directory")
		}
	})

	return za
}

func (a *ZipArchiver) ArchiveDirectory(dirPath, zipPath string) error {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	err := filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		relPath := info.Name()

		writer, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})

	if err != nil {
		zipWriter.Close()
		return err
	}

	if err := zipWriter.Close(); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(zipPath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	err = os.WriteFile(zipPath, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
