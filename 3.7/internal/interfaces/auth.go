package interfaces

import (
	"github.com/golang-jwt/jwt/v5"

	"warehouse-control/internal/domain"
)

type TokenProvide interface {
	Generate(user *domain.User) (string, error)
	Parse(token string) (*jwt.Token, error)
}
