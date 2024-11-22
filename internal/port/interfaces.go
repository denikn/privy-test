package port

import "mime/multipart"

type Storage interface {
	SaveFile(file multipart.File, filename string) (string, error)
	GetFilepath(filename string) string
}

type ImageService interface {
	SaveImage(file multipart.File, filename string) (string, error)
	ProcessImage(filename string) (string, int, error)
	GetFacesCount(filename string) (int, error)
	GetProcessedImagePath(filename string) string
}
