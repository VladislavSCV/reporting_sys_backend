package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RoleMiddleware проверяет, имеет ли пользователь нужную роль для доступа к маршруту
func RoleMiddleware(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем roleID из контекста
		roleID, exists := c.Get("roleID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found in context"})
			c.Abort()
			return
		}

		// Преобразуем roleID в строку (название роли)
		var roleValue string
		switch roleID {
		case 1:
			roleValue = "admin"
		case 2:
			roleValue = "teacher"
		case 3:
			roleValue = "student"
		default:
			c.JSON(http.StatusForbidden, gin.H{"error": "Unknown role"})
			c.Abort()
			return
		}

		// Проверяем, есть ли роль пользователя в списке разрешенных ролей
		allowed := false
		for _, allowedRole := range allowedRoles {
			if roleValue == allowedRole {
				allowed = true
				break
			}
		}

		// Если роль не разрешена, возвращаем ошибку
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		// Передаем управление следующему обработчику
		c.Next()
	}
}
