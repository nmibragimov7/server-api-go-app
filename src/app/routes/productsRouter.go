package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmibragimov7/go-app-server/src/app/controllers"
)

func ProductsRoutes(incoming *gin.RouterGroup) {
	incoming.Use()
	incoming.GET("/products", controllers.GetProducts)
	incoming.GET("/products/:id", controllers.GetProduct)
	incoming.POST("/products", controllers.AddProduct)
}
