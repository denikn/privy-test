package application

import (
	"io/ioutil"
	"os"
	"mime/multipart"
)

// MockStorage is a mock implementation of port.Storage for testing purposes.
type MockStorage struct{}

func (m *MockStorage) SaveFile(file multipart.File, filename string) (string, error) {
	// Simple implementation to simulate file storage
	tempDir := os.TempDir()
	tempFile := tempDir + "/" + filename
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(tempFile, data, 0644)
	if err != nil {
		return "", err
	}
	return tempFile, nil
}

func (m *MockStorage) GetFilepath(filename string) string {
	return "/mocked/file/path/" + filename
}
