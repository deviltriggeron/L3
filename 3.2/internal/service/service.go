package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"
	"net/url"
	e "shortener/internal/entity"
)

type ShortenerService struct {
	db *sql.DB
}

func NewShortenerService(db *sql.DB) ShortenerService {
	return ShortenerService{
		db: db,
	}
}

func (s *ShortenerService) NewShorten(ctx context.Context, baseUrl string) (string, error) {
	_, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	alias := generateShortUrl()

	_, err = s.db.ExecContext(ctx, "INSERT INTO url(alias, original_url) VALUES ($1, $2)", alias, baseUrl)
	if err != nil {
		return "", err
	}

	return alias, nil
}

func (s *ShortenerService) Redirect(ctx context.Context, shortUrl string) (string, error) {
	var originalURL string

	err := s.db.QueryRow("SELECT original_url FROM url WHERE alias = $1", shortUrl).Scan(&originalURL)
	if err == sql.ErrNoRows {
		fmt.Println("url not found")
		return "", fmt.Errorf("url not found")
	} else if err != nil {
		log.Fatal(err)
		return "", err
	}

	return originalURL, nil
}

func (s *ShortenerService) LogTransition(shortUrl string, userAgent string, ip string) error {
	var urlID string
	err := s.db.QueryRow("SELECT id FROM url WHERE alias=$1", shortUrl).Scan(&urlID)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
        INSERT INTO analytics (url_id, user_agent, ip_address)
        VALUES ($1, $2, $3)
    `, urlID, userAgent, ip)

	return err
}

func (s *ShortenerService) GetAnalyticsData(ctx context.Context, shortUrl string) ([]e.Transition, error) {
	var transitions []e.Transition
	rows, err := s.db.Query(`
        SELECT *
        FROM analytics
        WHERE url_id = (SELECT id FROM url WHERE alias=$1)
        ORDER BY time_transitions DESC
    `, shortUrl)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t e.Transition
		if err := rows.Scan(&t.ID, &t.URLID, &t.UserAgent, &t.IPAddress, &t.TimeTransitions); err != nil {
			return nil, err
		}
		transitions = append(transitions, t)
	}
	return transitions, nil
}

func generateShortUrl() string {
	n := rand.IntN(10) + 4
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	code := make([]rune, n)
	for i := range code {
		code[i] = letters[rand.IntN(len(letters))]
	}
	return string(code)
}
