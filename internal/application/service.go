package application

import (
	"beprivytest/internal/port"
	"path/filepath"
	"mime/multipart"
	"os/exec"
	"strings"
	"fmt"
	"bufio"
	"errors"
	"os"
)

type imageService struct {
	storage port.Storage
}

func NewImageService(storage port.Storage) port.ImageService {
	return &imageService{storage: storage}
}

func (s *imageService) SaveImage(file multipart.File, filename string) (string, error) {
	originalPath, err := s.storage.SaveFile(file, filename)
	if err != nil {
		return "", err
	}

	return originalPath, nil
}

func (s *imageService) ProcessImage(filename string) (string, int, error) {
	originalPath := s.storage.GetFilepath(filename)
	originalDir := filepath.Dir(originalPath)
	processedFilename := "processed_" + filepath.Base(originalPath)
	processedPath := filepath.Join(originalDir, processedFilename)

	// Check if the processed image already exists
	if _, err := os.Stat(processedPath); err == nil {
		// Try to read face count from cache
		faceCount, err := s.getFaceCountFromCache(processedFilename)
		if err != nil {
			return processedPath, faceCount, nil
		}
		// If there's an error reading the cache, proceed to reprocess the image
		fmt.Println("Error reading cache: ", err)
	}

	cmd := exec.Command("python3", "scripts/face_detection.py", originalPath, processedPath)

	// Pipe to capture the output of a Python script
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", 0, err
	}
	
	if err := cmd.Start(); err != nil {
		return "", 0, err
	}

	// Reading output from python pipe
	scanner := bufio.NewScanner(stdout)
	var faceCount int
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Number of detected faces:") {
			// Parse face count from python
			_, err := fmt.Sscanf(line, "Number of detected faces: %d", &faceCount)
			if err != nil {
				return "", 0, err
			}
			break
		}
	}

	// waiting process until done
	if err := cmd.Wait(); err != nil {
		return "", 0, err
	}
	// Handle error if no face detection
	if faceCount == 0 {
		return "", 0, errors.New("no faces detected")
	}

	// Store face count in cache
	err = s.storeFaceCountInCache(processedFilename, faceCount)
	if err != nil {
		return "", 0, err
	}

	return processedPath, faceCount, nil
}

// Function to get the number of faces from the file cache
func (s *imageService) getFaceCountFromCache(filename string) (int, error) {
	cacheFilePath := s.getCacheFilePath(filename)
	data, err := os.ReadFile(cacheFilePath)
	if err != nil {
		return 0, err
	}
	
	var faceCount int
	_, err = fmt.Sscanf(string(data), "%d", &faceCount)
	if err != nil {
		return 0, err
	}
	
	return faceCount, nil
}

// Function to save face count to file cache
func (s *imageService) storeFaceCountInCache(filename string, faceCount int) error {
	cacheFilePath := s.getCacheFilePath(filename)
	data := fmt.Sprintf("%d", faceCount)
	err := os.WriteFile(cacheFilePath, []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}

// Function to get cache file path based on file name
func (s *imageService) getCacheFilePath(filename string) string {
	cacheDir := "cache"
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err = os.MkdirAll(cacheDir, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating cache directory: %v\n", err)
			return ""
		}
	}
	return filepath.Join(cacheDir, filename+".cache")
}

func (s *imageService) GetProcessedImagePath(filename string) string {
	return s.storage.GetFilepath("processed_" + filename)
}

// New method to get face count for a given image
func (s *imageService) GetFacesCount(filename string) (int, error) {
	// Check if the processed image exists
	processedPath := s.GetProcessedImagePath(filename)
	if _, err := os.Stat(processedPath); err == nil {
		// Try to read face count from cache
		faceCount, err := s.getFaceCountFromCache("processed_" + filename)
		if err == nil {
			return faceCount, nil
		}
		return 0, err
	} else {
		// If processed image does not exist, process the image
		_, faceCount, err := s.ProcessImage(filename)
		if err != nil {
			return 0, err
		}
		return faceCount, nil
	}
}