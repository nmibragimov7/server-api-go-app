package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nmibragimov7/go-app-server/src/app/db"
	"github.com/nmibragimov7/go-app-server/src/app/middleware"
	"github.com/nmibragimov7/go-app-server/src/app/models"
	"github.com/nmibragimov7/go-app-server/src/app/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strings"
	"time"
)

var UsersCollection = db.OpenCollection(db.OpenDatabase(db.Client, "test"), "users")

func GetProfile(c *gin.Context) {
	token := c.GetHeader(middleware.AuthorizationHeaderKey)
	fields := strings.Fields(token)
	id := fmt.Sprintf("%v", service.Parse(fields[1]))
	fmt.Println(id)

	objectId, _ := primitive.ObjectIDFromHex(id)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var user models.User
	err := UsersCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user)

	fmt.Println(user)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь найден"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"profile": user})
}

func GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := UsersCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(ctx)

	var users []bson.M
	if err := cursor.All(ctx, &users); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		objectId, _ := primitive.ObjectIDFromHex(id)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var user models.User
		err := UsersCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь найден"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func SignIn(c *gin.Context) {
	var body models.User

	if err := c.BindJSON(&body); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var user models.User
	err := UsersCollection.FindOne(ctx, bson.M{"username": body.Username}).Decode(&user)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Неверный логин и пароль"})
		return
	}

	if user.ComparePassword(body.Password) {
		var tokenExpiresAt = time.Hour
		var refreshExpiresAt = time.Hour * 3

		token, err := service.JwtCreate(user.ID.Hex(), tokenExpiresAt)
		refresh, err := service.JwtCreate(user.ID.Hex(), refreshExpiresAt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(token)
		c.JSON(http.StatusOK, gin.H{
			"token":           token,
			"tokenExpireAt":   time.Now().Add(tokenExpiresAt),
			"refresh":         refresh,
			"refreshExpireAt": time.Now().Add(refreshExpiresAt),
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный логин и пароль"})
	}
}

func SingUp(c *gin.Context) {
	var body models.User

	if err := c.BindJSON(&body); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Неверный тип данных",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	encrypt, err := models.EncryptPassword(body.Password)
	if err != nil {
		log.Fatal(err)
	}

	user, err := UsersCollection.InsertOne(ctx, bson.D{
		{"first_name", body.FirstName},
		{"last_name", body.LastName},
		{"username", body.Username},
		{"password", encrypt},
		{"role", "user"},
		{"createdAt", time.Now()},
		{"updatedAt", time.Now()},
	})

	if err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusForbidden, err)
	}
	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"message": "Пользователь успешно зарегистрирован!",
	})
}
