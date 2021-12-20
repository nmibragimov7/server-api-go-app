package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nmibragimov7/go-app-server/src/app/db"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"time"
)

func GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client := db.ConnectDB()
	defer client.Disconnect(ctx)

	database := client.Database("test")
	fmt.Println("Database connect successful: ")

	collection := database.Collection("users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(ctx)

	var users []bson.M
	if err := cursor.All(ctx, &users); err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusForbidden, err)
	}

	c.IndentedJSON(http.StatusOK, users)
}
