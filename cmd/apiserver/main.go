package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nmibragimov7/go-app-server/src/app/controllers"
	"github.com/nmibragimov7/go-app-server/src/app/db"
)

type Config struct {
	port        string
	databaseUrl string
}

func NewConfig() *Config {
	return &Config{
		port:        ":8080",
		databaseUrl: "localhost:27017",
	}
}

func main() {
	config := NewConfig()
	db.ConnectDB(config.databaseUrl)
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	router.Use(cors.New(corsConfig))

	v1 := router.Group("/api")
	{
		v1.GET("/products", controllers.GetProducts)
		v1.POST("/products", controllers.PostProduct)
		v1.GET("/users", controllers.GetUsers)
	}

	router.Run("localhost" + config.port)
}
