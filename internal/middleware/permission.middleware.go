package middleware

import (
	"movie-app-go/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissionsInterface, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("No permissions in token"))
			c.Abort()
			return
		}

		permissions, ok := permissionsInterface.([]interface{})
		if !ok {
			c.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Invalid permissions format in token"))
			c.Abort()
			return
		}

		for _, p := range permissions {
			if p.(string) == permission {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, utils.ForbiddenResponse("Insufficient permissions: "+permission))
		c.Abort()
	}
}
