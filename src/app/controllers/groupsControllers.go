package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nmibragimov7/go-app-server/src/app/db"
	"github.com/nmibragimov7/go-app-server/src/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

var groupsCollection = db.OpenCollection(db.OpenDatabase(db.Client, "test"), "groups")

func GetGroups(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := groupsCollection.Find(ctx, bson.M{})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	defer cursor.Close(ctx)

	var groups []bson.M
	if err = cursor.All(ctx, &groups); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

func AddGroup(c *gin.Context) {
	var body models.Group

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

	_, err := groupsCollection.InsertOne(ctx, bson.D{
		{"name", body.Name},
		{"hash", body.Hash},
	})

	if err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusForbidden, err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Группа успешно добавлена!"})
}

func DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := groupsCollection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Произошла ошибка при удалении!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Группа успешно удалена!"})
}
