package storage

import (
	"beprivytest/internal/port"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type fileStorage struct {
	baseDir string
}

func NewFileStorage(baseDir string) port.Storage {
	return &fileStorage{baseDir: baseDir}
}

func (s *fileStorage) SaveFile(file multipart.File, filename string) (string, error) {
	destPath := filepath.Join(s.baseDir, filename)
	outFile, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		return "", err
	}

	return destPath, nil
}

func (s *fileStorage) GetFilepath(filename string) string {
	return filepath.Join(s.baseDir, filename)
}
