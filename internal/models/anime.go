package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AnimeDraft struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"id_user" json:"id_user"`
	AnimeID []string `bson:"id_anime" json:"id_anime"`
}