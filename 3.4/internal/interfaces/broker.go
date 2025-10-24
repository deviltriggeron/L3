package interfaces

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type EventPublisher interface {
	Consume(ctx context.Context, handler func(key string, value []byte) error)
	Produce(topic, key string, value []byte) error
	GetReader() *kafka.Reader
}
