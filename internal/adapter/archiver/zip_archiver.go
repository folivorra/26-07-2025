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
	zipDir string
	logger *slog.Logger
}

var _ Archiver = (*ZipArchiver)(nil)

func NewZipArchiver(a *app.App, zipDir string, logger *slog.Logger) *ZipArchiver {
	za := &ZipArchiver{
		a:      a,
		zipDir: zipDir,
		logger: logger,
	}

	za.a.RegisterCleanup(func(ctx context.Context) {
		if err := os.RemoveAll(za.zipDir); err != nil {
			za.logger.Warn("failed to remove archives directory")
			return
		}
		za.logger.Info("removed archives directory")
	})

	return za
}

func (a *ZipArchiver) ArchiveDirectory(dirPath string) error {
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

	if err := os.MkdirAll(filepath.Dir(a.zipDir), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	err = os.WriteFile(a.zipDir, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
