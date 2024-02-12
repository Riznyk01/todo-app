package todo_app

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

type Config struct {
	AccessTokenTtl  time.Duration
	RefreshTokenTtl time.Duration
	SigningKey      string
}

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewConfig(log *logrus.Logger) (cfg *Config, err error) {
	var SignKey string
	AccessTtl := setTokenTTL(log, os.Getenv("ACCESS_TTL"), "30m")
	RefreshTtl := setTokenTTL(log, os.Getenv("REFRESH_TTL"), "720h")

	if os.Getenv("SIGNING_KEY_FILE") != "" {
		SignKey, err = readSecret(os.Getenv("SIGNING_KEY_FILE"))
		if err != nil {
			return &Config{}, err
		}
	} else {
		SignKey = os.Getenv("SIGNING_KEY")
	}
	return &Config{
		AccessTokenTtl:  AccessTtl,
		RefreshTokenTtl: RefreshTtl,
		SigningKey:      SignKey,
	}, nil
}

func setTokenTTL(log *logrus.Logger, t, def string) time.Duration {
	TokenTTL, err := time.ParseDuration(t)
	if err != nil {
		TokenTTL, _ = time.ParseDuration(def)
		log.Errorf("Failed to set access/refresh token expiration duration. Using default value - %v", TokenTTL)
		return TokenTTL
	}
	return TokenTTL
}

func NewConfigPostgres() (cfg PostgresConfig, err error) {
	var dbPass string
	if os.Getenv("DB_PASSWORD_FILE") != "" {
		dbPass, err = readSecret(os.Getenv("DB_PASSWORD_FILE"))
		if err != nil {
			return PostgresConfig{}, err
		}
	} else {
		dbPass = os.Getenv("DB_PASSWORD")
	}
	return PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: dbPass,
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}, nil
}

func readSecret(path string) (secret string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	data := make([]byte, 64)
	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		secret = string(data[:n])
	}
	return secret, nil
}
