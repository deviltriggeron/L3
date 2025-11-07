package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"warehouse-control/internal/domain"
	"warehouse-control/internal/interfaces"
)

type jwtProvider struct {
	secret []byte
}

func NewJWTProvider(secret []byte) interfaces.TokenProvide {
	return &jwtProvider{
		secret: secret,
	}
}

func (p *jwtProvider) Generate(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secret)
}

func (p *jwtProvider) Parse(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return p.secret, nil
	})
}
