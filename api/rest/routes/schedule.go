package routes

import (
	"database/sql"
	"github.com/VladislavSCV/api/middleware"
	"github.com/VladislavSCV/api/rest/handlers"
	"github.com/gin-gonic/gin"
)

func SetupScheduleRoutes(router *gin.Engine, db *sql.DB) {
	scheduleGroup := router.Group("/api/schedule")
	{
		// Применение rate limiting к маршрутам
		scheduleGroup.Use(middleware.RateLimiterMiddleware())
		scheduleGroup.GET("/", handlers.GetSchedules(db))
		scheduleGroup.GET("/:id", handlers.GetScheduleByID(db))

		// Ограничение доступа для создания, обновления и удаления расписания только для администраторов
		scheduleGroup.POST("/", middleware.AuthMiddleware(db), middleware.RoleMiddleware([]string{"admin"}), handlers.CreateSchedule(db))
		scheduleGroup.PUT("/:id", middleware.AuthMiddleware(db), middleware.RoleMiddleware([]string{"admin"}), handlers.UpdateSchedule(db))
		scheduleGroup.DELETE("/:id", middleware.AuthMiddleware(db), middleware.RoleMiddleware([]string{"admin"}), handlers.DeleteSchedule(db))
	}
}
