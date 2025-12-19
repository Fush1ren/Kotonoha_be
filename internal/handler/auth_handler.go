package handler

import (
	"Kotonoha_be/internal/config"
	"Kotonoha_be/internal/models"
	"Kotonoha_be/internal/repository"
	"Kotonoha_be/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	userRepo    *repository.UserRepository
	refreshRepo repository.RefreshTokenRepo
	animeRepo *repository.AnimeRepository
}

func NewAuthHandler(u *repository.UserRepository, r repository.RefreshTokenRepo, a *repository.AnimeRepository) *AuthHandler {
	return &AuthHandler{u, r, a}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req struct {
		Name string `json:"name"`
		Username string `json:"username"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := utils.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := models.User{
		ID: primitive.NewObjectID(),
		Name: req.Name,
		Username: req.Username,
		Email: req.Email,
		Password: hash,
	}

	if err := h.userRepo.Create(ctx, &user); 
	err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.animeRepo.CreateAnimeDraft(ctx, user.ID); 
	err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create anime draft"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "register success"})
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req struct {
		Identity string `json:"identity" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	ctx.ShouldBindJSON(&req)

	user, err := h.userRepo.FindByUsernameOrEmail(ctx, req.Identity)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, _ := utils.GenerateAccessToken(user.ID)
	refreshToken, _ := utils.GenerateRefreshToken(user.ID)


	h.refreshRepo.Save(ctx, user.ID, refreshToken, config.RefreshTokenTTL)

	ctx.SetCookie("access_token", accessToken, 900, "/", "", true, true)
	ctx.SetCookie("refresh_token", refreshToken, int(config.RefreshTokenTTL.Seconds()), "/", "", true, true)

	ctx.JSON(200, gin.H{"message": "login success"})
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	rt, _ := ctx.Cookie("refresh_token")
	h.refreshRepo.Revoke(ctx, rt)

	ctx.SetCookie("access_token", "", -1, "/", "", true, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "", true, true)

	ctx.JSON(200, gin.H{"message": "logout success"})
}