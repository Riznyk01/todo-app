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
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

// TODO returning an error
func (s *AuthService) generatePasswordHash(pass string) string {
	fc := "Auth generatePasswordHash"

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("%s%s", fc, err)
	}
	return fmt.Sprintf("%s", passHash)
}
