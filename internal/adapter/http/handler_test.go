package http

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockImageService struct {
	mock.Mock
}

func (m *MockImageService) SaveImage(file multipart.File, filename string) (string, error) {
	args := m.Called(file, filename)
	return args.String(0), args.Error(1)
}

func (m *MockImageService) ProcessImage(filename string) (string, int, error) {
	args := m.Called(filename)
	return args.String(0), args.Int(1), args.Error(2)
}

func (m *MockImageService) GetProcessedImagePath(filename string) string {
	args := m.Called(filename)
	return args.String(0)
}

func (m *MockImageService) GetFacesCount(filename string) (int, error) {
	args := m.Called(filename)
	return args.Int(0), args.Error(1)
}

func TestUploadImage(t *testing.T) {
	mockService := new(MockImageService)
	handler := NewImageHandler(mockService)

	fileContent := "mock image content"
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.png")
	part.Write([]byte(fileContent))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	mockService.On("SaveImage", mock.Anything, "test.png").Return("", nil)

	handler.UploadImage(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertCalled(t, "SaveImage", mock.Anything, "test.png")
}

func TestGetFacesCount(t *testing.T) {
	mockService := new(MockImageService)
	handler := NewImageHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/faces?filename=test.png", nil)
	w := httptest.NewRecorder()

	mockService.On("ProcessImage", "test.png").Return("processed/test.png", 5, nil)

	handler.GetFacesCount(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	expected := `{"code":200,"result":{"processedPath":"processed/test.png","faceCount":5}}`
	actual := strings.TrimSpace(w.Body.String())
	assert.JSONEq(t, expected, actual)
}

func TestServeProcessedImage(t *testing.T) {
	mockService := new(MockImageService)
	handler := NewImageHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/processed-image?filename=test.png", nil)
	w := httptest.NewRecorder()

	mockService.On("GetProcessedImagePath", "test.png").Return("/path/to/processed/test.png")

	handler.ServeProcessedImage(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	mockService.AssertCalled(t, "GetProcessedImagePath", "test.png")
}
