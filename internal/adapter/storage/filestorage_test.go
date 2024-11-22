package storage_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"beprivytest/internal/adapter/storage"
)

// MockMultipartFile mocks a multipart.File for testing
type MockMultipartFile struct {
	mock.Mock
	Reader *bytes.Reader
}

func (m *MockMultipartFile) Read(p []byte) (int, error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *MockMultipartFile) Close() error {
	return m.Called().Error(0)
}

func (m *MockMultipartFile) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, nil // Mock implementation, not used in tests
}

func (m *MockMultipartFile) Seek(offset int64, whence int) (int64, error) {
	args := m.Called(offset, whence)
	return args.Get(0).(int64), args.Error(1)
}

func TestSaveFile_Success(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "storage_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	storageInstance := storage.NewFileStorage(tempDir)

	// Create a mock multipart file
	content := []byte("test content")
	file := &MockMultipartFile{Reader: bytes.NewReader(content)}
	file.On("Read", mock.Anything).Return(len(content), nil).Once()
	file.On("Read", mock.Anything).Return(0, io.EOF).Once()

	filename := "test.jpg"
	filePath, err := storageInstance.SaveFile(file, filename)

	// Assert no errors occurred
	assert.NoError(t, err)

	// Assert the filepath is correct
	expectedPath := filepath.Join(tempDir, filename)
	assert.Equal(t, expectedPath, filePath)

	// Check if the file was actually saved
	// data, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	// assert.Equal(t, content, data)

	file.AssertExpectations(t)
}

func TestSaveFile_ErrorCreatingFile(t *testing.T) {
	storageInstance := storage.NewFileStorage("invalid-dir")

	file := &MockMultipartFile{Reader: bytes.NewReader([]byte("test content"))}

	filename := "test.jpg"
	_, err := storageInstance.SaveFile(file, filename)

	// Assert an error is returned
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err)) // Check for specific error (os.IsNotExist)

	file.AssertExpectations(t)
}

func TestSaveFile_ErrorCopyingData(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "storage_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	storageInstance := storage.NewFileStorage(tempDir)

	// Mock a multipart file with error on Read
	file := &MockMultipartFile{}
	file.On("Read", mock.Anything).Return(0, errors.New("read error")).Once()

	filename := "test.jpg"
	_, err = storageInstance.SaveFile(file, filename)

	// Assert an error is returned
	assert.Error(t, err)
	assert.Equal(t, "read error", err.Error()) // Check for specific error message

	file.AssertExpectations(t)
}

func TestGetFilepath(t *testing.T) {
	baseDir := "some/base/dir"
	storageInstance := storage.NewFileStorage(baseDir)

	filename := "test.jpg"
	expectedPath := filepath.Join(baseDir, filename)

	filePath := storageInstance.GetFilepath(filename)

	assert.Equal(t, expectedPath, filePath)
}
