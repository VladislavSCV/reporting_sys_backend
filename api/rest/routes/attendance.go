package routes

import (
	"database/sql"
	"github.com/VladislavSCV/api/middleware"
	"github.com/VladislavSCV/api/rest/handlers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func SetupAttendanceRoutes(router *gin.Engine, db *sql.DB, cache *redis.Client) {
	attendanceGroup := router.Group("/api/attendance")
	{
		// Применение rate limiting к маршрутам
		attendanceGroup.Use(middleware.RateLimiterMiddleware())

		// Получение посещаемости по ID студента (доступно всем)
		attendanceGroup.GET("/student/:id", handlers.GetAttendanceByStudentID(db))
		// Получение посещаемости по ID группы (доступно всем)
		attendanceGroup.GET("/group/:id", handlers.GetAttendanceByGroupID(db))
		// Создание записи о посещаемости (доступно только преподавателям и администраторам)
		attendanceGroup.POST("/",
			middleware.AuthMiddleware(db),
			middleware.RoleMiddleware([]string{"teacher", "admin"}),
			handlers.CreateAttendance(db),
		)
		// Обновление записи о посещаемости (доступно только преподавателям и администраторам)
		attendanceGroup.PUT("/:id",
			middleware.AuthMiddleware(db),
			middleware.RoleMiddleware([]string{"teacher", "admin"}),
			handlers.UpdateAttendance(db),
		)
	}
}
