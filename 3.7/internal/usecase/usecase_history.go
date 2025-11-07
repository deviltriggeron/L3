package usecase

import (
	"warehouse-control/internal/domain"
	"warehouse-control/internal/interfaces"
)

type historyService struct {
	storage interfaces.HistoryRepository
}

func NewHistoryService(storage interfaces.HistoryRepository) HistoryService {
	return &historyService{
		storage: storage,
	}
}

func (h *historyService) GetHistory() ([]domain.ItemHistory, error) {
	items, err := h.storage.GetHistory()
	if err != nil {
		return nil, err
	}

	return items, nil
}
