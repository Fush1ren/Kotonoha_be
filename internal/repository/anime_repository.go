package repository

import (
	"Kotonoha_be/internal/models"
	"context"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnimeRepository struct {
	col *qmgo.Collection
}

func NewAnimeRepository(col *qmgo.Collection) *AnimeRepository {
	return &AnimeRepository{col}
}
// func NewAnimeRepository(db *qmgo.Database) *AnimeRepository {
// 	return &AnimeRepository{
// 		col: db.Collection("anime_draft"),
// 	}
// }

func (r *AnimeRepository) CreateAnimeDraft(ctx context.Context, userID primitive.ObjectID) error {
	draft := models.AnimeDraft {
		UserID: userID,
		AnimeID: []string{},
	}

	_, err := r.col.InsertOne(ctx, &draft)

	return err
}

func EnsureDraftIndex(col *qmgo.Collection) error {
	ctx := context.Background()

	return col.CreateOneIndex(
		ctx,
		options.IndexModel{
			Key: []string{"id_user"},
		},
	)
}

func (r *AnimeRepository) AddToDraft(ctx context.Context, draftID primitive.ObjectID, animeID string) error {
	return r.col.UpdateOne(
		ctx,
		bson.M{"_id": draftID},
		bson.M{
			"$addToSet": bson.M{
				"id_anime": animeID,
			},
		},
	)
}

func (r *AnimeRepository) RemoveAnime(ctx context.Context, draftID primitive.ObjectID, animeID string) error {
	return r.col.UpdateOne(
		ctx,
		bson.M{"_id": draftID},
		bson.M{
			"&pull": bson.M{
				"id_anime": animeID,
			},
		},
	)
}
// func (r *AnimeRepository) UpdateAnimeDraft(ctx context.Context, draftID primitive.ObjectID, animeID []string) error {
// 	return r.col.UpdateOne(
// 		ctx,
// 		bson.M{"_id": draftID},
// 		bson.M{
// 			"$addToSet": bson.M{
// 				"id_anime": bson.M{
// 					"$each": animeID,
// 				},
// 			},
// 		},
// 	)
// }

func (r *AnimeRepository) Create(ctx context.Context, anime *models.AnimeDraft) error {
	_, err := r.col.InsertOne(ctx, anime)
	return err
}

func (r *AnimeRepository) ListByUser(ctx context.Context, userID primitive.ObjectID) ([]models.AnimeDraft, error) {
	var result []models.AnimeDraft
	err := r.col.Find(ctx, bson.M{"id_user": userID}).All(&result)
	return result, err
}