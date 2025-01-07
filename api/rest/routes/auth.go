package routes

import (
	"database/sql"
	"github.com/VladislavSCV/api/middleware"
	"github.com/VladislavSCV/api/rest/handlers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func SetupAuthRoutes(router *gin.Engine, db *sql.DB, cache *redis.Client) {
	// Группа маршрутов для аутентификации
	auth := router.Group("/api/auth")
	{
		// Применение rate limiting к маршрутам
		auth.Use(middleware.RateLimiterMiddleware())

		// Маршруты
		auth.POST("/login", handlers.Login(db, cache))
		auth.POST("/registration", handlers.Registration(db))
		auth.POST("/verify", handlers.Verify(db))
		auth.GET("/", middleware.AuthMiddleware(db), handlers.GetCurrentUser(db, cache))
	}
}
