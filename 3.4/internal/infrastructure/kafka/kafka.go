package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"

	"imageprocessor/internal/domain"
	"imageprocessor/internal/interfaces"
)

type kafkaBroker struct {
	writer *kafka.Writer
	reader *kafka.Reader
	topic  string
}

func NewKafkaBroker(cfg domain.ConfigBroker) interfaces.EventPublisher {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Broker),
		Balancer: &kafka.LeastBytes{},
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.Broker},
		GroupID: cfg.GroupID,
		Topic:   cfg.Topic,
	})
	return &kafkaBroker{writer: writer, reader: reader, topic: cfg.Topic}
}

func (k *kafkaBroker) Consume(ctx context.Context, handler func(key string, value []byte) error) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer stopping...")
			if err := k.reader.Close(); err != nil {
				log.Printf("error kafka consumer stopping: %v", err)
				return
			}
			return
		default:
		}

		msg, err := k.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Kafka consumer stopped due to context cancellation")
				return
			}
			log.Println("Kafka read error:", err)
			continue
		}

		if err := handler(string(msg.Key), msg.Value); err != nil {
			log.Println("Handler error:", err)
		}
	}
}

func (k *kafkaBroker) Produce(topic, key string, value []byte) error {
	return k.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	})
}

func (k *kafkaBroker) GetReader() *kafka.Reader {
	return k.reader
}
