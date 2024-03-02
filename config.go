package todo_app

import (
	"errors"
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
	cfg = &Config{}

	cfg.AccessTokenTtl = setTokenTTL(log, os.Getenv("ACCESS_TTL"), "30m")
	cfg.RefreshTokenTtl = setTokenTTL(log, os.Getenv("REFRESH_TTL"), "720h")

	if os.Getenv("SIGNING_KEY_FILE") != "" {
		cfg.SigningKey, err = readSecret(os.Getenv("SIGNING_KEY_FILE"))
		if err != nil {
			return nil, err
		}
	} else {
		cfg.SigningKey = os.Getenv("SIGNING_KEY")
	}

	if cfg.SigningKey == "" {
		return nil, errors.New("signing key for the JWT is empty")
	}

	return cfg, nil
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

func NewConfigPostgres() (postgresCfg *PostgresConfig, err error) {
	postgresCfg = &PostgresConfig{}

	if os.Getenv("DB_PASSWORD_FILE") != "" {
		postgresCfg.Password, err = readSecret(os.Getenv("DB_PASSWORD_FILE"))
		if err != nil {
			return nil, err
		}
	} else {
		postgresCfg.Password = os.Getenv("DB_PASSWORD")
	}

	postgresCfg.Host = os.Getenv("DB_HOST")
	postgresCfg.Port = os.Getenv("DB_PORT")
	postgresCfg.Username = os.Getenv("DB_USERNAME")
	postgresCfg.DBName = os.Getenv("DB_NAME")
	postgresCfg.SSLMode = os.Getenv("DB_SSLMODE")

	if postgresCfg.Host == "" || postgresCfg.Port == "" || postgresCfg.Username == "" || postgresCfg.DBName == "" || postgresCfg.SSLMode == "" || postgresCfg.Password == "" {
		return nil, errors.New("some fields in the PostgreSQL config are empty")
	}

	return postgresCfg, nil
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
