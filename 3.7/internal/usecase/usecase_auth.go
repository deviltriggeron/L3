package usecase

import (
	"errors"

	"warehouse-control/internal/interfaces"
)

type authService struct {
	storage interfaces.UserRepository
	tp      interfaces.TokenProvide
}

func NewAuthService(tp interfaces.TokenProvide, storage interfaces.UserRepository) AuthService {
	return &authService{
		storage: storage,
		tp:      tp,
	}
}

func (s *authService) Login(username string, password string) (string, error) {
	user, err := s.storage.GetByUsername(username)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("user not found")
	}

	if !checkPassword(password, user.Pass) {
		return "", errors.New("invalid credentials")
	}

	return s.tp.Generate(user)
}

func checkPassword(password string, userPassword string) bool {
	return password == userPassword
}
