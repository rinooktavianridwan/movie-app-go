package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissionsInterface, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - no permissions in token"})
			c.Abort()
			return
		}

		permissions, ok := permissionsInterface.([]interface{})
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid permissions format in token"})
			c.Abort()
			return
		}

		for _, p := range permissions {
			if p.(string) == permission {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":               "Insufficient permissions",
			"required_permission": permission,
		})
		c.Abort()
	}
}
