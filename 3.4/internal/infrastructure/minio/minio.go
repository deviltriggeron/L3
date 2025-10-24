package minio

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"imageprocessor/internal/domain"
	"imageprocessor/internal/interfaces"
)

type minioClient struct {
	client          *minio.Client
	bucketOriginal  string
	bucketProcessed string
}

func NewMinioRepo(cfg domain.MinioCfg) interfaces.MinioRepository {
	mc, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.SSL,
	})
	if err != nil {
		log.Fatalf("failed to init MinIO client: %v", err)
	}
	return &minioClient{
		client:          mc,
		bucketOriginal:  cfg.Bucket,
		bucketProcessed: cfg.BucketProcessed,
	}
}

func (m *minioClient) InitMinio() error {
	ctx := context.Background()
	for _, bucket := range []string{m.bucketOriginal, m.bucketProcessed} {
		exists, err := m.client.BucketExists(ctx, bucket)
		if err != nil {
			return fmt.Errorf("error create bucket: %v", err)
		}
		if !exists {
			err := m.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
			if err != nil {
				return fmt.Errorf("error create bucket: %v", err)
			}
		}
	}
	return nil
}

func (m *minioClient) Create(file domain.FileDataType, options ...domain.ImageProcessOption) (string, error) {
	ctx := context.Background()
	id := uuid.New().String()
	objectName := fmt.Sprintf("%s_%s", id, file.FileName)

	reader := bytes.NewReader(file.Data)
	_, err := m.client.PutObject(ctx, m.bucketOriginal, objectName, reader, int64(len(file.Data)), minio.PutObjectOptions{
		ContentType: file.ContentType,
	})
	if err != nil {
		return "", fmt.Errorf("put object: %w", err)
	}
	return objectName, nil
}

func (m *minioClient) Get(objectID, variant string) ([]byte, string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s/%s", variant, objectID)

	obj, err := m.client.GetObject(ctx, m.bucketProcessed, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", err
	}
	defer func() {
		if err := obj.Close(); err != nil {
			log.Printf("error close minio object close: %v", err)
			return
		}
	}()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, obj)
	if err != nil {
		return nil, "", err
	}

	contentType := "image/jpeg"
	if strings.HasSuffix(strings.ToLower(objectID), ".png") {
		contentType = "image/png"
	}

	return buf.Bytes(), contentType, nil
}

func (m *minioClient) Delete(objectID string) error {
	ctx := context.Background()

	var errs []error
	for _, b := range []string{m.bucketOriginal, m.bucketProcessed} {
		err := m.client.RemoveObject(ctx, b, objectID, minio.RemoveObjectOptions{})
		if err != nil {
			errs = append(errs, fmt.Errorf("delete from original bucket failed: %w", err))
		}
	}

	var variants = []string{"resized", "thumb", "watermarked"}
	for _, variant := range variants {
		key := fmt.Sprintf("%s/%s", variant, objectID)
		if err := m.client.RemoveObject(ctx, m.bucketProcessed, key, minio.RemoveObjectOptions{}); err != nil {
			if minio.ToErrorResponse(err).Code != "NoSuchKey" {
				errs = append(errs, fmt.Errorf("delete from processed/%s failed: %w", variant, err))
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("delete errors: %v", errs)
	}

	return nil
}

func (m *minioClient) ProcessImage(ctx context.Context, objectName string, options ...domain.ImageProcessOption) error {
	obj, err := m.client.GetObject(ctx, m.bucketOriginal, objectName, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("get object: %w", err)
	}
	defer func() {
		if err := obj.Close(); err != nil {
			log.Printf("error close minio object close: %v", err)
			return
		}
	}()

	img, format, err := image.Decode(obj)
	if err != nil {
		return fmt.Errorf("decode image: %w", err)
	}

	for _, option := range options {
		switch option {
		case domain.ProcessResize:
			resized := imaging.Resize(img, 1080, 0, imaging.Lanczos)
			if err := m.uploadProcessed(ctx, objectName, resized, format, "resized"); err != nil {
				return err
			}
		case domain.ProcessThumbnail:
			thumb := imaging.Thumbnail(img, 200, 200, imaging.Lanczos)
			if err := m.uploadProcessed(ctx, objectName, thumb, format, "thumbnail"); err != nil {
				return err
			}
		case domain.ProcessWatermark:
			watermarked := addWatermark(img)
			if err := m.uploadProcessed(ctx, objectName, watermarked, format, "watermarked"); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *minioClient) uploadProcessed(ctx context.Context, objectName string, img image.Image, format string, variant string) error {
	buf := new(bytes.Buffer)
	var contentType string

	switch strings.ToLower(format) {
	case "png":
		contentType = "image/png"
		if err := png.Encode(buf, img); err != nil {
			return fmt.Errorf("encode png: %w", err)
		}
	case "jpeg", "jpg":
		contentType = "image/jpeg"
		if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 90}); err != nil {
			return fmt.Errorf("encode jpeg: %w", err)
		}
	default:
		contentType = "image/jpeg"
		if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 90}); err != nil {
			return fmt.Errorf("encode default jpeg: %w", err)
		}
	}

	key := fmt.Sprintf("%s/%s", variant, objectName)
	_, err := m.client.PutObject(ctx, m.bucketProcessed, key, buf, int64(buf.Len()), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("put object failed: %w", err)
	}

	return nil
}

func addWatermark(img image.Image) image.Image {
	nrgba := imaging.Clone(img)
	bounds := nrgba.Bounds()
	bar := image.Rect(bounds.Min.X, bounds.Max.Y-40, bounds.Max.X, bounds.Max.Y)
	draw.Draw(nrgba, bar, &image.Uniform{C: color.RGBA{0, 0, 0, 100}}, image.Point{}, draw.Over)
	return nrgba
}
