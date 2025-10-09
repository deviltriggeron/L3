package service

import (
	"fmt"
	"log"
	e "notifier/internal/entity"
	msgbroker "notifier/internal/rabbitMQ"
	"strconv"
	"sync"
	"time"
)

const (
	maxRetries = 5
	baseDelay  = 5 * time.Second
)

type NotifierService struct {
	data   map[int]*e.Notification
	broker *msgbroker.Broker
	mu     sync.Mutex
	nextID int
}

func NewNotifierService() *NotifierService {
	b := msgbroker.Connect()

	s := &NotifierService{
		data:   make(map[int]*e.Notification),
		broker: b,
		mu:     sync.Mutex{},
		nextID: 1,
	}

	go s.worker()

	return s
}

func (s *NotifierService) NewNotification(notifyHandle e.NotifierHandle) *e.Notification {
	s.mu.Lock()
	defer s.mu.Unlock()

	notify := e.Notification{
		ID:       s.nextID,
		Message:  notifyHandle.Message,
		SendAt:   notifyHandle.SendAt,
		Status:   e.Pending,
		Attempts: 0,
	}
	s.nextID++

	n := notify
	s.data[n.ID] = &n

	s.broker.Produce(notify)
	return &n
}

func (s *NotifierService) GetStatus(stringID string) (string, error) {
	id, err := strconv.Atoi(stringID)
	if err != nil {
		log.Printf("cannot parse ID: %s", stringID)
		return "", fmt.Errorf("cannot parse ID: %s", stringID)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	gettingID, ok := s.data[id]
	if !ok {
		log.Printf("not found ID: %d", id)
		return "", fmt.Errorf("not found ID: %d", id)
	}

	return gettingID.Status, nil
}

func (s *NotifierService) DeleteNotify(stringID string) error {
	id, err := strconv.Atoi(stringID)
	if err != nil {
		log.Printf("cannot parse ID: %s", stringID)
		return fmt.Errorf("cannot parse ID: %s", stringID)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[id]
	if !ok {
		log.Printf("not found ID: %s", stringID)
		return fmt.Errorf("not found ID: %s", stringID)
	}
	s.data[id].Status = e.Cancelled
	return nil
}

func (s *NotifierService) worker() {
	for idStr := range s.broker.Consume() {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Invalid message ID: %s", idStr)
			continue
		}

		s.mu.Lock()
		n, ok := s.data[id]
		s.mu.Unlock()

		if !ok {
			log.Printf("Notification %d not found in memory", id)
			continue
		}

		go s.processNotification(n)
	}
}

func (s *NotifierService) processNotification(n *e.Notification) {
	sleepDuration := n.SendAt.Sub(time.Now())
	if sleepDuration > 0 {
		log.Printf("Notification %d sleeping for %v until %s", n.ID, sleepDuration, n.SendAt.Format(time.RFC3339))
		time.Sleep(sleepDuration)
	}

	for {
		if n.Status != e.Cancelled {
			err := sendNotification(n)
			if err != nil {
				n.Attempts++
				if n.Attempts > maxRetries {
					n.Status = e.Failed
					log.Printf("Notification %d failed after %d attempts", n.ID, n.Attempts)
					return
				}

				delay := time.Duration(1<<n.Attempts) * baseDelay
				log.Printf("Retry %d for notification %d after %v", n.Attempts, n.ID, delay)
				time.Sleep(delay)
				continue
			}

			n.Status = e.Sent
			log.Printf("Notification %d sent successfully", n.ID)
			return
		}
	}
}

func sendNotification(n *e.Notification) error {
	log.Printf("Sending notification %d: %s", n.ID, n.Message)
	return nil
}
