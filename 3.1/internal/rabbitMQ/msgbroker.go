package msgbroker

import (
	"log"
	e "notifier/internal/entity"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	amqpURL   = "amqp://guest:guest@localhost:5672/"
	queueName = "notifications"
)

type Broker struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func Connect() *Broker {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Error connect message broker: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel message broker: %v", err)
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	return &Broker{
		conn: conn,
		ch:   ch,
	}
}

func (b *Broker) Close() {
	if b.ch != nil {
		b.ch.Close()
	}

	if b.conn != nil {
		b.conn.Close()
	}
}

func (b *Broker) Produce(msg e.Notification) {
	id := strconv.Itoa(msg.ID)
	err := b.ch.Publish(
		"", queueName, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(id),
		})
	if err != nil {
		log.Fatalf("Failed to produce message: %v", err)
	}
}

func (b *Broker) Consume() chan string {
	msgs, err := b.ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	out := make(chan string)
	go func() {
		for msg := range msgs {
			out <- string(msg.Body)
		}
	}()
	return out
}
