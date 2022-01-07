package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmibragimov7/go-app-server/src/app/controllers"
	"github.com/nmibragimov7/go-app-server/src/app/middleware"
)

func GroupsRoutes(incoming *gin.RouterGroup) {
	routes := incoming.Group("").Use(middleware.AuthMiddleware())
	routes.GET("/groups", controllers.GetGroups)
	routes.POST("/groups", controllers.AddGroup)
	routes.DELETE("/groups/:id", controllers.DeleteGroup)
}
