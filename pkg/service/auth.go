package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"time"
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
	passHash, err := generatePasswordHash(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = passHash
	return s.repo.CreateUser(user)
}
func (s *AuthService) GenerateToken(username, password string) (string, error) {
	var (
		ErrInvalidCredentials = errors.New("invalid credentials")
	)

	fc := "GenerateToken"
	var tokenTtl time.Duration
	tokenHours, err := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	if err != nil {
		log.Errorf("%s failed to set token expiration duration\n %s. Using default value - 12 Hours.", fc, err)
		tokenHours = 12
	}
	tokenTtl = time.Duration(tokenHours) * time.Hour

	user, err := s.repo.GetUser(username)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Errorf("%s invalid credentials: %v", fc, err)
		return "", fmt.Errorf("%s: %w", fc, ErrInvalidCredentials)
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(tokenTtl).Unix()
	claims["ess"] = time.Now().Unix()
	claims["userid"] = user.Id

	tokenString, err := token.SignedString([]byte(os.Getenv("SIGNING_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func generatePasswordHash(pass string) (string, error) {
	fc := "Auth generatePasswordHash"

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("%s failed to generate password hash: %v", fc, err)
		return "", fmt.Errorf("failed to generate password hash %w", err) //err
	}
	return fmt.Sprintf("%s", passHash), nil
}
