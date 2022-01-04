package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName string             `json:"first_name" validate="required"`
	LastName  string             `json:"last_name" validate="required"`
	Username  string             `json:"username" validate="required"`
	Password  string             `json:"password" validate="required, min=0"`
	Token     string             `json:"token"`
	Role      string             `json:"role" validate="required, eq=admin|eq=user"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func EncryptPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (user *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
}
