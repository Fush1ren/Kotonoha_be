package config

import (
	"os"
	"time"
)

var (
	AccessTokenSecret  string
	RefreshTokenSecret string

	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

func MustLoadJWTConfig() {
	AccessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
	RefreshTokenSecret = os.Getenv("REFRESH_TOKEN_SECRET")

	if AccessTokenSecret == "" || RefreshTokenSecret == "" {
		panic("JWT secret is missing in environment")
	}
}
