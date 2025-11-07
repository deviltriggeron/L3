package postgresql

import (
	"context"
	"database/sql"
	"errors"

	"warehouse-control/internal/domain"
	"warehouse-control/internal/interfaces"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) interfaces.UserRepository {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) GetByUsername(username string) (*domain.User, error) {
	var user domain.User

	q := `
		SELECT *
		FROM users
		WHERE username = $1
	`

	err := u.db.QueryRowContext(context.Background(), q, username).Scan(&user.ID, &user.Username, &user.Pass, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
