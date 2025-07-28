package archiver

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
)

type ZipArchiver struct{}

func NewZipArchiver() *ZipArchiver {
	return &ZipArchiver{}
}

func (a *ZipArchiver) ArchiveDirectory(dirPath, zipPath string) error {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	err := filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Пропускаем директории
		if info.IsDir() {
			return nil
		}

		// Открываем файл
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Получаем относительный путь, чтобы структура директорий сохранялась
		relPath, err := filepath.Rel(dirPath, filePath)
		if err != nil {
			return err
		}

		// Создаём файл внутри архива
		writer, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// Копируем содержимое файла в архив
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

	err = os.WriteFile(zipPath, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
