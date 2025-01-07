package middleware

import (
	"database/sql"
	"github.com/VladislavSCV/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		// Валидируем токен
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Сохраняем userID и roleID в контексте
		c.Set("userID", claims.UserID)
		c.Set("roleID", claims.RoleID)

		// Передаем управление следующему обработчику
		c.Next()
	}
}
