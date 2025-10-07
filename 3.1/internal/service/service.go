package service

import (
	"fmt"
	"log"
	e "notifier/internal/entity"
	"strconv"
	"sync"
)

type NotifierService struct {
	data map[int]*e.Notification
	mu   sync.Mutex
	next int
}

func NewNotifierService() *NotifierService {
	return &NotifierService{
		data: make(map[int]*e.Notification),
		mu:   sync.Mutex{},
		next: 1,
	}
}

func (s *NotifierService) NewNotification(notify e.Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[s.next] = &notify
}

func (s *NotifierService) GetStatus(stringID string) (string, error) {
	id, err := strconv.Atoi(stringID)
	if err != nil {
		log.Printf("cannot parse id: %s", stringID)
		return "", fmt.Errorf("cannot parse id: %s", stringID)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	gettingID, ok := s.data[id]
	if !ok {
		log.Printf("not found id: %d", id)
		return "", fmt.Errorf("not found id: %d", id)
	}

	return gettingID.Status, nil
}

func (s *NotifierService) DeleteNotify(stringID string) error {
	id, err := strconv.Atoi(stringID)
	if err != nil {
		log.Printf("cannot parse id: %s", stringID)
		return fmt.Errorf("cannot parse id: %s", stringID)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, id)

	return nil
}
