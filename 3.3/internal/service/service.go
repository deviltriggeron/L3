package service

import (
	e "commentTree/internal/entity"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type CommentsService struct {
	db *sql.DB
}

func NewCommentsService(db *sql.DB) *CommentsService {
	return &CommentsService{
		db: db,
	}
}

func (s *CommentsService) Comments(ctx context.Context, c e.CommentResponse) (*e.Comments, error) {
	if c.UserName == "" || c.Comment == "" {
		return nil, fmt.Errorf("not suitable username or comment")
	}

	var err error
	comments := &e.Comments{
		UserName: c.UserName,
		Comment:  c.Comment,
	}
	if c.ParentID == uuid.Nil {
		q := `
			INSERT INTO comments (user_name, comment, path)
			VALUES ($1, $2, REPLACE(gen_random_uuid()::text, '-', '')::ltree)
			RETURNING comment_id, path;
		`
		err = s.db.QueryRowContext(ctx, q, c.UserName, c.Comment).
			Scan(&comments.CommentID, &comments.Path)
	} else {
		q := `
			WITH parent AS (SELECT path FROM comments WHERE comment_id = $1)
			INSERT INTO comments (user_name, comment, parent_id, path)
			SELECT $2, $3, $1, parent.path || REPLACE(gen_random_uuid()::text, '-', '')::ltree
			FROM parent
			RETURNING comment_id, parent_id;
		`
		err = s.db.QueryRowContext(ctx, q, c.ParentID, c.UserName, c.Comment).
			Scan(&comments.CommentID, &comments.ParentID)
	}

	if err != nil {
		return nil, err
	}

	return comments, err
}

func (s *CommentsService) GetComments(ctx context.Context, stringID string) ([]e.Comments, error) {
	id, err := uuid.Parse(stringID)
	if err != nil {
		return nil, fmt.Errorf("invalid comment ID: %w", err)
	}

	q := `
		SELECT *
		FROM comments
		WHERE path <@ (SELECT path FROM comments WHERE comment_id = $1)
		  AND comment_id != $1
		ORDER BY path;
	`

	rows, err := s.db.QueryContext(ctx, q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []e.Comments
	for rows.Next() {
		var c e.Comments
		if err := rows.Scan(&c.CommentID, &c.ParentID, &c.UserName, &c.Comment, &c.Path, &c.Date); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return comments, nil
}

func (s *CommentsService) DeleteComments(ctx context.Context, stringID string) error {
	id, err := uuid.Parse(stringID)
	if err != nil {
		return err
	}
	q := `
		DELETE FROM comments WHERE comment_id = $1;
	`
	_, err = s.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *CommentsService) GetAllParentComments(ctx context.Context) ([]e.Comments, error) {
	q := `
		SELECT *
		FROM comments
		WHERE parent_id IS NULL
		ORDER BY date;
	`

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []e.Comments
	for rows.Next() {
		var c e.Comments
		err = rows.Scan(&c.CommentID, &c.ParentID, &c.UserName, &c.Comment, &c.Path, &c.Date)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
