package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmibragimov7/go-app-server/src/app/controllers"
)

func ProductsRoutes(incoming *gin.RouterGroup) {
	incoming.GET("/products", controllers.GetProducts)
	incoming.POST("/products", controllers.PostProduct)
}
