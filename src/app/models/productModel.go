package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID    primitive.ObjectID `bson:"_id"`
	Title string             `json:"title" validate="required"`
	Price int64              `json:"price" validate="required, min=0"`
	Group string             `json:"group" validate="required"`
}
