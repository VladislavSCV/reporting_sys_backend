package routes

import (
	"database/sql"
	"github.com/VladislavSCV/api/middleware"
	"github.com/VladislavSCV/api/rest/handlers"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, db *sql.DB) {
	userGroup := router.Group("/api/user")
	{
		// Применение rate limiting к маршрутам
		userGroup.Use(middleware.RateLimiterMiddleware())
		userGroup.GET("/", handlers.GetUsers(db))
		userGroup.GET("/students", handlers.GetStudents(db))
		userGroup.GET("/teachers", handlers.GetTeachers(db))
		userGroup.GET("/:id", handlers.GetUserByID(db))

		// Ограничение доступа для обновления и удаления пользователей только для администраторов
		userGroup.PUT("/:id", middleware.AuthMiddleware(db), middleware.RoleMiddleware([]string{"admin"}), handlers.UpdateUser(db))
		userGroup.DELETE("/:id", middleware.AuthMiddleware(db), middleware.RoleMiddleware([]string{"admin"}), handlers.DeleteUser(db))
	}
}
