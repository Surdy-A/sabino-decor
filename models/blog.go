package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Tittle      string             `json:"title" bson:"title"`
	Content     string             `json:"content" bson:"content"`
	PublishedAt time.Time          `json:"published_at" bson:"published_at"`
	Image       string             `json:"image" bson:"image"`
	Comments    []string           `json:"comments" bson:"comments"`
	Likes       []string           `json:"likes" bson:"likes"`
}
