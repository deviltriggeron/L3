package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"sales-tracker/internal/domain"
)

type storage struct {
	db *sql.DB
}

func InitDB(cfg domain.DBConfig) (domain.StorageRepositoryExtended, error) {
	var err error

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.DB,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Database initialized successfully!")
	return &storage{db: db}, nil
}

func (s *storage) Insert(item domain.Item) error {
	ctx := context.Background()
	q := `
		INSERT INTO item (
			id, type, amount, category, description, create_date
		)
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	_, err := s.db.ExecContext(ctx, q, item.ID, item.Type, item.Amount, item.Category, item.Description, item.Date)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) Select(itemID uuid.UUID) (*domain.Item, error) {
	var item domain.Item

	ctx := context.Background()
	q := `
		SELECT *
		FROM item
		WHERE id = $1
	`

	err := s.db.QueryRowContext(ctx, q, itemID).Scan(&item.ID, &item.Type, &item.Amount, &item.Category, &item.Description, &item.Date)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *storage) Update(itemID uuid.UUID, item domain.Item) error {
	ctx := context.Background()
	q := `
		UPDATE item SET
			type = $1, amount = $2, category = $3, description = $4, create_date = $5
		WHERE id = $6
	`

	_, err := s.db.ExecContext(ctx, q, item.Type, item.Amount, item.Category, item.Description, item.Date, itemID)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) Delete(itemID uuid.UUID) error {
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

func (s *storage) GetAll(filter domain.ItemFilter) ([]domain.Item, error) {
	ctx := context.Background()

	q, args := createQueryForItem(filter)
	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error close row in postgresql: %v", err)
		}
	}()

	var items []domain.Item
	for rows.Next() {
		var it domain.Item
		var t string
		if err := rows.Scan(&it.ID, &t, &it.Amount, &it.Category, &it.Description, &it.Date); err != nil {
			return nil, err
		}
		it.Type = domain.ItemType(t)
		items = append(items, it)
	}
	return items, nil
}

func createQueryForItem(filter domain.ItemFilter) (string, []interface{}) {
	args := []interface{}{}

	q := `
		SELECT *
		FROM item
		WHERE 1 = 1
	`
	idx := 1

	if filter.From != "" {
		q += fmt.Sprintf(" AND create_date >= $%d", idx)
		args = append(args, filter.From)
		idx++
	}
	if filter.To != "" {
		q += fmt.Sprintf(" AND create_date <= $%d", idx)
		args = append(args, filter.To)
		idx++
	}
	if filter.Type != "" {
		q += fmt.Sprintf(" AND type = $%d", idx)
		args = append(args, string(filter.Type))
		idx++
	}

	if filter.Category != "" {
		q += fmt.Sprintf(" AND category = $%d", idx)
		args = append(args, filter.Category)
		idx++
	}
	q += " ORDER BY create_date DESC"
	if filter.Limit > 0 {
		q += fmt.Sprintf(" LIMIT $%d", idx)
		args = append(args, filter.Limit)
		idx++
	}

	if filter.Offset > 0 {
		q += fmt.Sprintf(" OFFSET $%d", idx)
		args = append(args, filter.Offset)
	}

	return q, args
}

func (s *storage) GetAnalytics(filter domain.AnalyticsFilter) (*domain.AnalyticsResult, error) {
	var res domain.AnalyticsResult
	ctx := context.Background()
	q := `
		SELECT
		COALESCE(SUM(amount)::double precision,0) as sum,
		COALESCE(AVG(amount)::double precision,0) as avg,
		COALESCE(COUNT(*)::int,0) as count,
		COALESCE(percentile_cont(0.5) WITHIN GROUP (ORDER BY amount)::double precision,0) as median,
		COALESCE(percentile_cont(0.9) WITHIN GROUP (ORDER BY amount)::double precision,0) as p90
		FROM item
		WHERE 1=1
	`
	args := []interface{}{}
	idx := 1

	if filter.From != "" {
		q += fmt.Sprintf(" AND create_date >= $%d", idx)
		args = append(args, filter.From)
		idx++
	}
	if filter.To != "" {
		q += fmt.Sprintf(" AND create_date <= $%d", idx)
		args = append(args, filter.To)
		idx++
	}
	if filter.Category != "" {
		q += fmt.Sprintf(" AND category = $%d", idx)
		args = append(args, filter.Category)
		idx++
	}
	if filter.Type != "" {
		q += fmt.Sprintf(" AND type = $%d", idx)
		args = append(args, filter.Type)
	}

	row := s.db.QueryRowContext(ctx, q, args...)

	if err := row.Scan(
		&res.Sum,
		&res.Avg,
		&res.Count,
		&res.Median,
		&res.Percentile90,
	); err != nil {
		return nil, err
	}
	return &res, nil
}
