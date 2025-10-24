package interfaces

import (
	"context"

	"imageprocessor/internal/domain"
)

type MinioRepository interface {
	InitMinio() error
	Create(file domain.FileDataType, options ...domain.ImageProcessOption) (string, error)
	Get(objectID, variant string) ([]byte, string, error)
	Delete(objectID string) error
	ProcessImage(ctx context.Context, objectName string, options ...domain.ImageProcessOption) error
}
