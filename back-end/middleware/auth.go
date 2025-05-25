package middleware

import (
	"RushOrder/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminAuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		admin, err := service.GetAdminSession(c, db)
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if admin == nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("admin", admin)
		c.Next()
	}
}
