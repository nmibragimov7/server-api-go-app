package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nmibragimov7/go-app-server/src/app/service"
	"net/http"
	"strings"
)

const (
	AuthorizationHeaderKey = "Authorization"
	AuthorizationType      = "bearer"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(AuthorizationHeaderKey)
		if len(token) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован!"})
			c.Abort()
			return
		}

		fields := strings.Fields(token)
		if len(fields) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный формат токена!"})
			c.Abort()
			return
		}

		tokenType := strings.ToLower(fields[0])
		if tokenType != AuthorizationType {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неподдерживаемый тип токена!"})
			c.Abort()
			return
		}

		access := service.Verify(fields[1])
		if access == false {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Отказано в доступе!"})
			c.Abort()
			return
		}

		c.Next()
	}
}
