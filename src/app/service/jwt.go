package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const secretKey = "secret"

func JwtCreate(userId string, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(duration).Unix(),
		Subject:   userId,
	})

	return token.SignedString([]byte(secretKey))
}

func Verify(incomingToken string) bool {
	token, err := jwt.Parse(incomingToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err == nil && token.Valid {
		return true
	} else {
		return false
	}
}

func Parse(incomingToken string) interface{} {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(incomingToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		fmt.Println(err)
	}

	for key, val := range claims {
		if key == "sub" {
			id := val
			return id
		}
	}

	return nil
}
