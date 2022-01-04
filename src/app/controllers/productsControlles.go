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

func GetProducts(c *gin.Context) {
	group := c.Query("group")

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
		log.Fatal(err)
	}

	defer cursor.Close(ctx)

	var products []bson.M
	if err := cursor.All(ctx, &products); err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusForbidden, err)
	}

	c.IndentedJSON(http.StatusOK, products)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Товар найден"})
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

	product, err := productsCollection.InsertOne(ctx, bson.D{
		{"title", body.Title},
		{"price", body.Price},
		{"group", body.Group},
	})

	if err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusForbidden, err)
	}
	c.JSON(http.StatusOK, gin.H{"product": product})
}
