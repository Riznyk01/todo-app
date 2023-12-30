package service

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	todoapp "todo-app"
	"todo-app/pkg/repository"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todoapp.User) (int, error) {
	passHash, err := s.generatePasswordHash(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = passHash
	return s.repo.CreateUser(user)
}

func (s *AuthService) generatePasswordHash(pass string) (string, error) {
	fc := "Auth generatePasswordHash"

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("%s failed to generate password hash\n %s", fc, err)
		return "", err
	}
	return fmt.Sprintf("%s", passHash), nil
}
