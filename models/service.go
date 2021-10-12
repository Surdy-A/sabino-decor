package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Image       string             `json:"image" bson:"image"`
}
