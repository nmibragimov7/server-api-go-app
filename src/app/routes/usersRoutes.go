package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmibragimov7/go-app-server/src/app/controllers"
)

func UsersRoutes(incoming *gin.RouterGroup) {
	incoming.GET("/users", controllers.GetUsers)
}
