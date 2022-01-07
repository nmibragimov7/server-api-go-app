package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nmibragimov7/go-app-server/src/app/db"
	"github.com/nmibragimov7/go-app-server/src/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

var productsCollection = db.OpenCollection(db.OpenDatabase(db.Client, "test"), "products")
var validate = validator.New()

func FetchGroups() []models.Group {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := groupsCollection.Find(ctx, bson.M{})

	if err != nil {
		fmt.Println(err)
	}

	defer cursor.Close(ctx)

	var groups []models.Group
	if err = cursor.All(ctx, &groups); err != nil {
		fmt.Println(err)
	}

	return groups
}

func GetCounts(c *gin.Context) {
	groups := FetchGroups()
	type Item struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	var counts []Item

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := productsCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	defer cursor.Close(ctx)

	var products []bson.M
	if err = cursor.All(ctx, &products); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	counts = append(counts, Item{"all", len(products)})

	for _, group := range groups {
		cursor, err = productsCollection.Find(ctx, bson.M{"group": group.Name})
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		if err = cursor.All(ctx, &products); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		counts = append(counts, Item{group.Hash, len(products)})
	}

	c.JSON(http.StatusOK, gin.H{"counts": counts})
}

func GetProducts(c *gin.Context) {
	group := c.Query("group")
	//type pagination struct {
	//	limit       int
	//	pages       int
	//	currentPage int
	//	count       int
	//}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var cursor *mongo.Cursor
	var err error = nil

	if group != "" {
		cursor, err = productsCollection.Find(ctx, bson.M{"group": group})
	} else {
		cursor, err = productsCollection.Find(ctx, bson.M{})
	}
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	defer cursor.Close(ctx)

	var products []bson.M
	if err = cursor.All(ctx, &products); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var product models.Product
	err := productsCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&product)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": product})
}

func AddProduct(c *gin.Context) {
	var body models.Product

	if err := c.BindJSON(&body); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Неверный тип данных",
		})
		return
	}

	//err := validate.Struct(body)
	//if err != nil {
	//	log.Fatal(err)
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := productsCollection.InsertOne(ctx, bson.D{
		{"title", body.Title},
		{"price", body.Price},
		{"group", body.Group},
	})

	if err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusForbidden, err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Товар успешно добавлен!"})
}

func EditProduct(c *gin.Context) {
	var body models.Product

	if err := c.BindJSON(&body); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Неверный тип данных!",
		})
		return
	}

	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := productsCollection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.D{{"$set", bson.D{
		{"title", body.Title},
		{"price", body.Price},
		{"group", body.Group},
	}}})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Произошла ошибка при обновлении!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Товар обновлен успешно!"})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := productsCollection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Произошла ошибка при удалении!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Товар успешно удален!"})
}
