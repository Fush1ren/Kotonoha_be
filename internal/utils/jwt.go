package utils

import (
	"Kotonoha_be/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	accessSecret  = []byte("ACCESS_SECRET_KEY")
	refreshSecret = []byte("REFRESH_SECRET_KEY")
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID primitive.ObjectID) (string, error) {
	return generateToken(
		userID,
		config.AccessTokenTTL,
		config.AccessTokenSecret,
	)
}

func GenerateRefreshToken(userID primitive.ObjectID) (string, error) {
	return generateToken(
		userID,
		config.RefreshTokenTTL,
		config.RefreshTokenSecret,
	)
}

func generateToken(
	userID primitive.ObjectID,
	ttl time.Duration,
	secret string,
) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.Hex(),
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(secret))
}