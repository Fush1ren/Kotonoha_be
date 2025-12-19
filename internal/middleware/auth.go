package middleware

import (
	"Kotonoha_be/internal/config"
	"Kotonoha_be/internal/repository"
	"Kotonoha_be/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(refreshRepo repository.RefreshTokenRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// =============================
		// Use Access Token
		// =============================
		accessToken, err := ctx.Cookie("access_token")
		if err == nil {
			token, err := jwt.Parse(accessToken, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(config.AccessTokenSecret), nil
			})

			if err == nil && token.Valid {
				claims := token.Claims.(jwt.MapClaims)
				ctx.Set("id_user", claims["sub"])
				ctx.Next()
				return
			}
		}

		// =============================
		// Access token invalid / expired and refresh
		// =============================
		refreshToken, err := ctx.Cookie("refresh_token")
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{
				"error": "unauthorized",
			})
			return
		}

		rtData, err := refreshRepo.Validate(ctx, refreshToken)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{
				"error": "refresh token invalid",
			})
			return
		}

		// =============================
		// Generate new access token
		// =============================
		newAccessToken, err := utils.GenerateAccessToken(rtData.UserID)
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{
				"error": "failed to generate access token",
			})
			return
		}

		ctx.SetCookie(
			"access_token",
			newAccessToken,
			int(config.AccessTokenTTL.Seconds()),
			"/",
			"",
			true,
			true,
		)

		ctx.Set("id_user", rtData.UserID.Hex())
		ctx.Next()
	}
}