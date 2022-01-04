package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmibragimov7/go-app-server/src/app/controllers"
	"github.com/nmibragimov7/go-app-server/src/app/middleware"
)

func UsersRoutes(incoming *gin.RouterGroup) {
	routes := incoming.Group("/").Use(middleware.AuthMiddleware())
	routes.GET("me", controllers.GetProfile)

	incoming.GET("/users", controllers.GetUsers)
	incoming.GET("/users/:id", controllers.GetUser())
}
