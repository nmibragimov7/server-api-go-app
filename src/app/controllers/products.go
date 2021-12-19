package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+"localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("test")
	fmt.Println("Database connect successful: ")

	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)

	collection := db.Collection("products")

	var cursor *mongo.Cursor

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

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+"localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("test")
	fmt.Println("Database connect successful: ")

	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)

	collection := db.Collection("products")

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
