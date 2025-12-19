package main

import (
	"Kotonoha_be/internal/config"
	"Kotonoha_be/internal/config/database"
	"Kotonoha_be/internal/handler"
	"Kotonoha_be/internal/repository"
	"Kotonoha_be/internal/router"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if errnv := godotenv.Load(); errnv != nil {
		log.Println("no .env file found, using system env")
	}
	config.MustLoadJWTConfig()

	mongoURI := os.Getenv("MONGO_URI")
	mongoDB := os.Getenv("MONGO_DB")
	port := os.Getenv("APP_PORT")
	db, errDB := database.Connect(mongoURI, mongoDB)

	if errDB != nil {
		log.Fatalf("failed to connect mongodb: %v", errDB)
	}

	r := gin.Default()
	
	userRepo := repository.NewUserRepository(db.Collection("users"))
	refreshRepo := repository.NewRefreshTokenRepository(db.Collection("refresh_tokens"))
	animeRepo := repository.NewAnimeRepository(db.Collection("anime_draft"))
	
	animeHandler := handler.NewAnimeHandler(animeRepo)
	authHandler := handler.NewAuthHandler(userRepo, refreshRepo, animeRepo)

	if err := repository.EnsureDraftIndex(
		db.Collection("anime_draft"),
	); err != nil {
		log.Fatal("failed to create anime draft index:", err)
	}

	router.AnimeRoutes(r, animeHandler, refreshRepo)
	router.AuthRoutes(r, authHandler)

	r.Run(":" + port)
}