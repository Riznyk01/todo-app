package todo_app

import (
	"time"
)

type TokenConfig struct {
	AccessTokenTtl  time.Duration
	RefreshTokenTtl time.Duration
}

func NewConfig() *TokenConfig {
	return &TokenConfig{}
}
