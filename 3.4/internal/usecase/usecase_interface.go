package usecase

import (
	"imageprocessor/internal/domain"
)

type ImageProcService interface {
	Upload(file domain.FileDataType, options ...domain.ImageProcessOption) (string, error)
	Get(id, variant string) ([]byte, string, error)
	DeleteImage(id string) error
}
