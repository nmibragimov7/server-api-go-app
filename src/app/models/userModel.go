package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName string             `json:"firstName" validate="required"`
	LastName  string             `json:"lastName" validate="required"`
	Username  string             `json:"username" validate="required"`
	Password  string             `json:"password" validate="required"`
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

func (user *User) BeforeSend() User {
	copy := User{user.ID, user.FirstName, user.LastName, user.Username, "", user.Role, user.CreatedAt, user.UpdatedAt}
	return copy
}
