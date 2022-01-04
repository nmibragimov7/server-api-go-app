package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmibragimov7/go-app-server/src/app/controllers"
)

func AuthRoutes(incoming *gin.RouterGroup) {
	incoming.POST("/sign-in", controllers.SignIn)
	incoming.POST("/sign-up", controllers.SingUp)
}
