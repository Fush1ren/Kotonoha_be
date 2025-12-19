package handler

import (
	"Kotonoha_be/internal/models"
	"Kotonoha_be/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnimeHandler struct {
	repo *repository.AnimeRepository
}

func NewAnimeHandler(repo *repository.AnimeRepository) *AnimeHandler {
	return &AnimeHandler{repo: repo}
}

func (h *AnimeHandler) CreateDraft(ctx *gin.Context) {
	var req struct {
		UserID  string `json:"id_user"`
		AnimeID []string `json:"id_anime"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	anime := models.AnimeDraft{
		UserID:  uid,
		AnimeID: req.AnimeID,
	}

	if err := h.repo.Create(ctx, &anime); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, anime)
}

func (h *AnimeHandler) AddAnimeToDraft(ctx *gin.Context) {
	draftID, _ := primitive.ObjectIDFromHex(ctx.Param("id_draft"))

	var req struct {
		AnimeID string `json:"id_anime"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.AddToDraft(ctx, draftID, req.AnimeID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (h *AnimeHandler) GetDraft(ctx *gin.Context) {
	// Ambil user_id dari context (string)
	userIDStr := ctx.GetString("id_user")

	// Convert ke ObjectID
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	draft, err := h.repo.ListByUser(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "draft not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, draft)
}

func (h *AnimeHandler) DeleteAnimeFromDraft(ctx *gin.Context) {
	draftID, err := primitive.ObjectIDFromHex(ctx.Param("id_draft"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	animeID := ctx.Param("id_anime")
	if animeID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "anime_id required"})
		return
	}

	if err := h.repo.RemoveAnime(ctx, draftID, animeID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "anime removed from draft"})
}