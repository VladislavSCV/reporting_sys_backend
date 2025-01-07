package routes

import (
	"database/sql"
	"github.com/VladislavSCV/api/middleware"
	"github.com/VladislavSCV/api/rest/handlers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func SetupGradeRoutes(router *gin.Engine, db *sql.DB, cache *redis.Client) {
	gradeGroup := router.Group("/api/grades")
	{
		// Применение rate limiting к маршрутам
		gradeGroup.Use(middleware.RateLimiterMiddleware())
		gradeGroup.GET("/student/:id", handlers.GetGradesByStudentID(db))
		gradeGroup.GET("/group/:id", handlers.GetGradesByGroupID(db))

		// Ограничение доступа для создания, обновления и удаления оценок только для преподавателей
		gradeGroup.POST("/", middleware.AuthMiddleware(db), middleware.RoleMiddleware([]string{"teacher"}), handlers.CreateGrade(db))
		gradeGroup.PUT("/:id", middleware.AuthMiddleware(db), middleware.RoleMiddleware([]string{"teacher"}), handlers.UpdateGrade(db))
		gradeGroup.DELETE("/:id", middleware.AuthMiddleware(db), middleware.RoleMiddleware([]string{"teacher"}), handlers.DeleteGrade(db))
	}
}
