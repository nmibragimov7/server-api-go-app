package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmibragimov7/go-app-server/src/app/controllers"
	"github.com/nmibragimov7/go-app-server/src/app/middleware"
)

func AuthRoutes(incoming *gin.RouterGroup) {
	routes := incoming.Group("").Use(middleware.AuthMiddleware())
	routes.GET("/me", controllers.GetProfile)
	routes.PUT("/me", controllers.EditProfile)
	routes.PUT("/change-password", controllers.ChangePassword)

	incoming.POST("/sign-in", controllers.SignIn)
	incoming.POST("/sign-up", controllers.SingUp)
	incoming.GET("/refresh", controllers.Refresh)
}
