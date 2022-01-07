package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmibragimov7/go-app-server/src/app/controllers"
	"github.com/nmibragimov7/go-app-server/src/app/middleware"
)

func UsersRoutes(incoming *gin.RouterGroup) {
	routes := incoming.Group("").Use(middleware.AuthMiddleware())
	routes.GET("/me", controllers.GetProfile)
	routes.GET("/users", controllers.GetUsers)
	routes.PUT("/users/:id", controllers.EditUser)
	routes.DELETE("/users/:id", controllers.DeleteUser)
	routes.GET("/users/:id", controllers.GetUser())
}
