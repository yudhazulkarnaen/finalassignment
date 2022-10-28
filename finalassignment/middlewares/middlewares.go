package middlewares

import (
	"errors"
	"net/http"

	"finalassignment.id/finalassignment/utils/token"
	"github.com/gin-gonic/gin"
)

const errorMessageStr = "error_message"

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		if err != nil {
			status := http.StatusUnauthorized
			if errors.Is(err, token.ErrNoToken) {
				status = http.StatusBadRequest
			}
			c.AbortWithStatusJSON(status, gin.H{
				errorMessageStr: err.Error(),
			})
			return
		}
		c.Next()
	}
}
