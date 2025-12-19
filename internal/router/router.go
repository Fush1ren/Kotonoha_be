package router

import (
	"Kotonoha_be/internal/handler"
	"Kotonoha_be/internal/middleware"
	"Kotonoha_be/internal/repository"

	"github.com/gin-gonic/gin"
)

func AnimeRoutes(r *gin.Engine, animeHandler *handler.AnimeHandler, refreshRepo repository.RefreshTokenRepo,) {
	anime := r.Group("/api/draft")
	anime.Use(middleware.AuthMiddleware(refreshRepo))
	{
		anime.GET("", animeHandler.GetDraft)
		anime.PUT("/:id_draft/anime", animeHandler.AddAnimeToDraft)
		anime.DELETE("/:id_draft/anime/:id_anime", animeHandler.DeleteAnimeFromDraft)
	}
}

func AuthRoutes(r *gin.Engine, authHandler *handler.AuthHandler) {
	auth := r.Group("/api/auth") 
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/logout", authHandler.Logout)
	}
}