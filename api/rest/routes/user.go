package routes

import (
	"database/sql"
	"github.com/VladislavSCV/api/middleware"
	"github.com/VladislavSCV/api/rest/handlers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func SetupUserRoutes(router *gin.Engine, db *sql.DB, cache *redis.Client) {
	userGroup := router.Group("/api/user")
	{
		// Применение rate limiting к маршрутам
		userGroup.Use(middleware.RateLimiterMiddleware())
		userGroup.GET("/", handlers.GetUsers(db, cache))
		userGroup.GET("/students", handlers.GetStudents(db, cache))
		userGroup.GET("/teachers", handlers.GetTeachers(db, cache))
		userGroup.GET("/:id", handlers.GetUserByID(db, cache))

		// Ограничение доступа для обновления и удаления пользователей только для администраторов
		userGroup.PUT("/:id", middleware.AuthMiddleware(db), middleware.RoleMiddleware([]string{"admin"}), handlers.UpdateUser(db, cache))
		userGroup.DELETE("/:id", middleware.AuthMiddleware(db), middleware.RoleMiddleware([]string{"admin"}), handlers.DeleteUser(db, cache))
	}
}
