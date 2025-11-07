package postgresql

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"

	"warehouse-control/internal/domain"
	"warehouse-control/internal/interfaces"
)

type itemRepo struct {
	db *sql.DB
}

func NewItemRepo(db *sql.DB) interfaces.ItemRepository {
	return &itemRepo{
		db: db,
	}
}

func (s *itemRepo) Select(itemID uuid.UUID) (*domain.Item, error) {
	var item domain.Item

	q := `
		SELECT *
		FROM item
		WHERE id = $1
	`

	err := s.db.QueryRowContext(context.Background(), q, itemID).Scan(&item.ID, &item.Product, &item.Price, &item.Description, &item.Count, &item.CreateDate)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *itemRepo) SelectAll() ([]domain.Item, error) {
	var items []domain.Item

	q := `
		SELECT *
		FROM item
	`

	rows, err := s.db.QueryContext(context.Background(), q)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error close row in database: %v", err)
		}
	}()

	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(&item.ID, &item.Product, &item.Price, &item.Description, &item.Count, &item.CreateDate); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (s *itemRepo) Insert(item domain.Item) error {
	q := `
		INSERT INTO item (
			id, product, price, description, count, create_date
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := s.db.ExecContext(context.Background(), q, item.ID, item.Product, item.Price, item.Description, item.Count, item.CreateDate)
	if err != nil {
		return err
	}

	return nil
}

func (s *itemRepo) Update(itemID uuid.UUID, item domain.Item) error {
	q := `
		UPDATE item SET 
			product = $1, price = $2, description = $3, count = $4, create_date = $5
		WHERE id = $6
	`

	_, err := s.db.ExecContext(context.Background(), q, item.Product, item.Price, item.Description, item.Count, item.CreateDate, itemID)
	if err != nil {
		return err
	}

	return nil
}
func (s *itemRepo) Delete(itemID uuid.UUID) error {
	ctx := context.Background()
	q := `
		DELETE 
		FROM item
		WHERE id = $1
	`

	_, err := s.db.ExecContext(ctx, q, itemID)
	if err != nil {
		return err
	}

	return nil
}
