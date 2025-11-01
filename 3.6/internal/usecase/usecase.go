package usecase

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"sales-tracker/internal/domain"
)

type trackService struct {
	storage domain.StorageRepositoryExtended
}

func NewTrackService(db domain.StorageRepositoryExtended) TrackerServiceExtended {
	return &trackService{
		storage: db,
	}
}

func (s *trackService) Insert(item domain.Item) (uuid.UUID, error) {
	id := uuid.New()
	item.ID = id
	if item.Amount <= 0 {
		return uuid.Nil, fmt.Errorf("the quantity cannot be less than 1")
	}
	err := s.storage.Insert(item)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *trackService) Get(stringID string) (*domain.Item, error) {
	itemID, err := uuid.Parse(stringID)
	if err != nil {
		return nil, err
	}

	item, err := s.storage.Select(itemID)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *trackService) Update(stringID string, item domain.Item) error {
	itemID, err := uuid.Parse(stringID)
	if err != nil {
		return err
	}

	err = s.storage.Update(itemID, item)
	if err != nil {
		return err
	}

	return nil
}

func (s *trackService) Delete(stringID string) error {
	itemID, err := uuid.Parse(stringID)
	if err != nil {
		return err
	}

	err = s.storage.Delete(itemID)
	if err != nil {
		return err
	}

	return nil
}

func (s *trackService) GetAll(from string, to string, category string, typ string, limit string, offset string) ([]domain.Item, error) {
	itemFilter := domain.ItemFilter{
		From:     from,
		To:       to,
		Category: category,
		Type:     domain.ItemType(typ),
	}

	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return nil, err
		}
		itemFilter.Limit = limitInt
	}
	if offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return nil, err
		}
		itemFilter.Offset = offsetInt
	}

	items, err := s.storage.GetAll(itemFilter)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *trackService) GetAnalytics(from string, to string, category string, typ string) (*domain.AnalyticsResult, error) {
	if from != "" {
		_, err := time.Parse("2006-01-02", from)
		if err != nil {
			return nil, err
		}
	}
	if to != "" {
		_, err := time.Parse("2006-01-02", to)
		if err != nil {
			return nil, err
		}
	}

	analyticsFilter := domain.AnalyticsFilter{
		From:     from,
		To:       to,
		Category: category,
		Type:     domain.ItemType(typ),
	}

	result, err := s.storage.GetAnalytics(analyticsFilter)
	if err != nil {
		return nil, err
	}

	return result, nil
}
