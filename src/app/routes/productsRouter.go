package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmibragimov7/go-app-server/src/app/controllers"
	"github.com/nmibragimov7/go-app-server/src/app/middleware"
)

func ProductsRoutes(incoming *gin.RouterGroup) {
	routes := incoming.Group("").Use(middleware.AuthMiddleware())
	routes.GET("/products", controllers.GetProducts)
	routes.GET("/products/:id", controllers.GetProduct)
	routes.POST("/products", controllers.AddProduct)
	routes.PUT("/products/:id", controllers.EditProduct)
	routes.DELETE("/products/:id", controllers.DeleteProduct)
	routes.GET("/products/count", controllers.GetCounts)
}
