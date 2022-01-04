package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nmibragimov7/go-app-server/src/app/db"
	"github.com/nmibragimov7/go-app-server/src/app/routes"
	"log"
	"os"
)

type Config struct {
	url string
}

func NewConfig() *Config {
	return &Config{
		url: "localhost:8080",
	}
}

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	BackendUrl := os.Getenv("BACKEND_URL")

	config := NewConfig()
	db.ConnectDB()
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"*"}
	corsConfig.AllowHeaders = []string{"*"}
	router.Use(cors.New(corsConfig))

	v1 := router.Group("/api")
	{
		routes.ProductsRoutes(v1)
		routes.UsersRoutes(v1)
		routes.AuthRoutes(v1)
	}

	if BackendUrl == "" {
		BackendUrl = config.url
	}

	router.Run(BackendUrl)
}
