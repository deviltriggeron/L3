package usecase

import (
	"encoding/json"
	"log"

	"imageprocessor/internal/domain"
	"imageprocessor/internal/interfaces"
)

type imageProcService struct {
	minioService   interfaces.MinioRepository
	eventPublisher interfaces.EventPublisher
}

func NewImageProcService(m interfaces.MinioRepository, e interfaces.EventPublisher) ImageProcService {
	return &imageProcService{
		minioService:   m,
		eventPublisher: e,
	}
}

func (s *imageProcService) Upload(file domain.FileDataType, options ...domain.ImageProcessOption) (string, error) {
	name, err := s.minioService.Create(file)
	if err != nil {
		return "", err
	}

	task := domain.ImageTask{
		ID:      name,
		Options: options,
	}
	data, _ := json.Marshal(task)

	if err := s.eventPublisher.Produce("image-tasks", name, data); err != nil {
		log.Printf("Kafka produce failed: %v", err)
	}

	return name, nil
}

func (s *imageProcService) Get(id, variant string) ([]byte, string, error) {
	return s.minioService.Get(id, variant)
}

func (s *imageProcService) DeleteImage(id string) error {
	err := s.minioService.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
