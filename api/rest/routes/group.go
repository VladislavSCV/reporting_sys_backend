package routes

import (
	"database/sql"
	"github.com/VladislavSCV/api/middleware"
	"github.com/VladislavSCV/api/rest/handlers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func SetupGroupRoutes(router *gin.Engine, db *sql.DB, cache *redis.Client) {
	groupGroup := router.Group("/api/group")
	{
		// Применение rate limiting к маршрутам
		groupGroup.Use(middleware.RateLimiterMiddleware())

		// Получение списка всех групп
		groupGroup.GET("/", handlers.GetGroups(db, cache))

		// Получение информации о группе (студенты и расписание)
		groupGroup.GET("/:id", handlers.GetGroupByID(db, cache))

		// Ограничение доступа для создания группы только для администраторов
		groupGroup.POST("/", middleware.AuthMiddleware(db), middleware.RoleMiddleware([]string{"admin"}), handlers.CreateGroup(db, cache))

		// Ограничение доступа для обновления группы только для администраторов
		groupGroup.PUT("/:id", middleware.AuthMiddleware(db), middleware.RoleMiddleware([]string{"admin"}), handlers.UpdateGroup(db, cache))

		// Ограничение доступа для удаления группы только для администраторов
		groupGroup.DELETE("/:id", middleware.AuthMiddleware(db), middleware.RoleMiddleware([]string{"admin"}), handlers.DeleteGroup(db, cache))
	}
}
