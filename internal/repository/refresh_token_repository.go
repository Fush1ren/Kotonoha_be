package repository

import (
	"Kotonoha_be/internal/models"
	"context"
	"time"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type RefreshTokenRepo interface {
	Save(ctx context.Context, userID primitive.ObjectID, token string, ttl time.Duration) error
	Validate(ctx context.Context, token string) (*models.RefreshToken, error)
	Revoke(ctx context.Context, token string) error
	RevokeAll(ctx context.Context, userID primitive.ObjectID) error
}

type RefreshTokenRepository struct {
	col *qmgo.Collection
}

func NewRefreshTokenRepository(col *qmgo.Collection) RefreshTokenRepo {
	return &RefreshTokenRepository{col}
}

func (r *RefreshTokenRepository) Save(ctx context.Context, userID primitive.ObjectID, token string, ttl time.Duration) error {
	_, err := r.col.InsertOne(ctx, models.RefreshToken{
		UserID:    userID,
		Token:     token,
		Revoked:   false,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
	})
	return err
}

func (r *RefreshTokenRepository) Validate(ctx context.Context, token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.col.Find(ctx, bson.M{
		"token": token,
		"revoked": false,
		"expires_at": bson.M{"$gt": time.Now()},
	}).One(&rt)
	return &rt, err
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	return r.col.UpdateOne(ctx,
		bson.M{"token": token},
		bson.M{"$set": bson.M{"revoked": true}},
	)
}

func (r *RefreshTokenRepository) RevokeAll(
	ctx context.Context,
	userID primitive.ObjectID,
) error {

	_, err := r.col.UpdateAll(
		ctx,
		bson.M{"id_user": userID},
		bson.M{"$set": bson.M{"revoked": true}},
	)

	return err
}
