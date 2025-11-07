package postgresql

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"warehouse-control/internal/domain"
	"warehouse-control/internal/interfaces"
)

type historyRepo struct {
	db *sql.DB
}

func NewHistoryRepo(db *sql.DB) interfaces.HistoryRepository {
	return &historyRepo{
		db: db,
	}
}

func (r *historyRepo) GetHistory() ([]domain.ItemHistory, error) {
	var history []domain.ItemHistory

	q := `
		SELECT *
		FROM items_history
	`

	rows, err := r.db.QueryContext(context.Background(), q)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error close row in database: %v", err)
		}
	}()

	for rows.Next() {
		var h domain.ItemHistory
		var oldJSON, newJSON []byte

		if err := rows.Scan(&h.ID, &h.ItemID, &h.Action, &oldJSON, &newJSON, &h.ChangedBy, &h.ChangedAt); err != nil {
			return nil, err
		}

		if len(oldJSON) > 0 {
			if err := json.Unmarshal(oldJSON, &h.OldData); err != nil {
				return nil, err
			}
		}

		if len(newJSON) > 0 {
			if err := json.Unmarshal(newJSON, &h.NewData); err != nil {
				return nil, err
			}
		}

		history = append(history, h)
	}

	return history, nil
}
