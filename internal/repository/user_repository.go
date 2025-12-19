package repository

import (
	"Kotonoha_be/internal/models"
	"context"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository struct {
	col *qmgo.Collection
}

func NewUserRepository(col *qmgo.Collection) *UserRepository {
	return &UserRepository{col}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.col.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.col.Find(ctx, bson.M{"username": username}).One(&user)
	return &user, err
}

func (r *UserRepository) FindByUsernameOrEmail(
	ctx context.Context,
	identity string,
) (*models.User, error) {
	var user models.User

	err := r.col.Find(ctx, bson.M{
		"$or": []bson.M{
			{"username": identity},
			{"email": identity},
		},
	}).One(&user)

	return &user, err
}