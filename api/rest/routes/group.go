package routes

import (
	"database/sql"
	"github.com/VladislavSCV/api/middleware"
	"github.com/VladislavSCV/api/rest/handlers"
	"github.com/gin-gonic/gin"
)

func SetupGroupRoutes(router *gin.Engine, db *sql.DB) {
	groupGroup := router.Group("/api/group")
	{
		// Применение rate limiting к маршрутам
		groupGroup.Use(middleware.RateLimiterMiddleware())
		groupGroup.GET("/", handlers.GetGroups(db))
		groupGroup.GET("/:id", handlers.GetGroupByID(db)) // получение информации о группе (студенты)

		// Ограничение доступа для создания группы только для администраторов
		groupGroup.POST("/", middleware.AuthMiddleware(db), middleware.RoleMiddleware([]string{"admin"}), handlers.CreateGroup(db))

		// Ограничение доступа для обновления и удаления группы только для администраторов
		groupGroup.PUT("/:id", middleware.AuthMiddleware(db), middleware.RoleMiddleware([]string{"admin"}), handlers.UpdateGroup(db))
		groupGroup.DELETE("/:id", middleware.AuthMiddleware(db), middleware.RoleMiddleware([]string{"admin"}), handlers.DeleteGroup(db))
	}
}
