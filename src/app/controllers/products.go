package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nmibragimov7/go-app-server/src/app/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
	"unicode/utf8"
)

type Product struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Price int64  `json:"price"`
	Group string `json:"group"`
}

func GetProducts(c *gin.Context) {
	//id := c.Param("id")
	group := c.Query("group")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client := db.ConnectDB()
	defer client.Disconnect(ctx)

	database := client.Database("test")
	fmt.Println("Database connect successful: ")

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := database.Collection("products")

	var cursor *mongo.Cursor
	var err error = nil

	if group != "" {
		cursor, err = collection.Find(ctx, bson.M{"group": group})
	} else {
		cursor, err = collection.Find(ctx, bson.M{})
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

func PostProduct(c *gin.Context) {
	var body Product

	if err := c.BindJSON(&body); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client := db.ConnectDB()
	defer client.Disconnect(ctx)

	database := client.Database("test")
	fmt.Println("Database connect successful: ")

	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)

	collection := database.Collection("products")

	product, err := collection.InsertOne(ctx, bson.D{
		{"title", body.Title},
		{"price", body.Price},
		{"group", body.Group},
	})
	fmt.Println(utf8.ValidString(body.Group))
	if err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusForbidden, err)
	}
	c.IndentedJSON(http.StatusOK, product)
}
