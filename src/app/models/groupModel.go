package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Group struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `json:"name" validate="required"`
	Hash string             `json:"hash" validate="required"`
}
