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
func (s *AuthService) ExistsUser(email string) (bool, error) {
	ex, err := s.repo.ExistsUser(email)
	if err != nil {
		return false, err
	}
	return ex, nil
}
func (s *AuthService) CheckUserPassword(email, password string) error {
	fc := "CheckUserPassword"

	user, err := s.repo.GetUser(email)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.log.Errorf("%s invalid credentials: %v", fc, err)
		return err
	}

	return nil
}
func (s *AuthService) GenerateTokenPair(email string) (string, string, error) {
	fc := "GenerateTokenPair"

	user, err := s.repo.GetUser(email)
	if err != nil {
		return "", "", err
	}

	accessMinutesStr, err := strconv.Atoi(os.Getenv("ACCESS_TTL"))
	if err != nil {
		s.log.Errorf("%s failed to set access token expiration duration\n %s. Using default value - 30 mins.", fc, err)
		accessMinutesStr = 30
	}
	accessTokenTtl := time.Duration(accessMinutesStr) * time.Minute

	refreshHoursStr, err := strconv.Atoi(os.Getenv("REFRESH_TTL"))
	if err != nil {
		s.log.Errorf("%s failed to set refresh token expiration duration\n %s. Using default value - 720 hours.", fc, err)
		refreshHoursStr = 720
	}
	refreshTokenTtl := time.Duration(refreshHoursStr) * time.Hour

	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["sub"] = user.Id
	accessClaims["iat"] = time.Now().Unix()
	accessClaims["exp"] = time.Now().Add(accessTokenTtl).Unix()

	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("SIGNING_KEY")))
	if err != nil {
		s.log.Errorf("%s error creating signed access token: %v", fc, err)
		return "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["sub"] = user.Id
	refreshClaims["mail"] = email
	refreshClaims["iat"] = time.Now().Unix()
	refreshClaims["exp"] = time.Now().Add(refreshTokenTtl).Unix()

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("SIGNING_KEY")))
	if err != nil {
		s.log.Errorf("%s error creating signed refresh token: %v", fc, err)
		return "", "", err
	}

	err = s.repo.UpdateRefreshTokenInDB(email, refreshTokenString)
	if err != nil {
		s.log.Errorf("%s error wrighting refresh token to DB: %v", fc, err)
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
func (s *AuthService) ParseToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("SIGNING_KEY")), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token structure")
	}
	useridInt, ok := claims["sub"].(float64)
	if !ok {
		return 0, errors.New("invalid token structure: userid is not a number")
	}
	return int(useridInt), nil
}
func (s *AuthService) CheckTokenInDB(refreshTokenString string) (string, error) {
	userMail, err := s.repo.CheckRefreshTokenInDB(refreshTokenString)
	if err != nil {
		return "", err
	}
	return userMail, nil
}
func generatePasswordHash(pass string) (string, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		//log.Errorf("%s failed to generate password hash: %v", fc, err)
		return "", fmt.Errorf("failed to generate password hash %w", err)
	}
	return fmt.Sprintf("%s", passHash), nil
}
