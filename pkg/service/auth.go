package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"time"
	todoapp "todo-app"
	"todo-app/pkg/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	log  *logrus.Logger
	repo repository.Authorization
}

func NewAuthService(log *logrus.Logger, repo repository.Authorization) *AuthService {
	return &AuthService{
		log:  log,
		repo: repo,
	}
}

func (s *AuthService) CreateUser(user todoapp.User) (int, error) {
	passHash, err := generatePasswordHash(user.Password)
	if err != nil {
		s.log.Errorf("%s failed to generate password hash: %v",
			"service. auth. generatePasswordHash", err)
		return 0, err
	}
	user.Password = passHash
	return s.repo.CreateUser(user)
}
func (s *AuthService) ExistsUser(username string) (bool, error) {
	ex, err := s.repo.UserExists(username)
	if err != nil {
		return false, err
	}
	return ex, nil
}
func (s *AuthService) GenerateToken(username, password string) (string, error) {
	fc := "GenerateToken"
	var tokenTtl time.Duration

	user, err := s.repo.GetUser(username)
	if err != nil {
		return "", err
	}

	tokenHours, err := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	if err != nil {
		s.log.Errorf("%s failed to set token expiration duration\n %s. Using default value - 12 Hours.", fc, err)
		tokenHours = 12
	}
	tokenTtl = time.Duration(tokenHours) * time.Hour

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.log.Errorf("%s invalid credentials: %v", fc, err)
		return "", ErrInvalidCredentials
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userid"] = user.Id
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(tokenTtl).Unix()
	claims["ess"] = time.Now().Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("SIGNING_KEY")))
	if err != nil {
		s.log.Errorf("%s error creating signed token: %v", fc, err)
		return "", err
	}
	return tokenString, nil
}
func generatePasswordHash(pass string) (string, error) {
	//fc := "Auth generatePasswordHash"
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		//log.Errorf("%s failed to generate password hash: %v", fc, err)
		return "", fmt.Errorf("failed to generate password hash %w", err)
	}
	return fmt.Sprintf("%s", passHash), nil
}
