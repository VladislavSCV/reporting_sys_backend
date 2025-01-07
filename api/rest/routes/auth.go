package routes

import (
	"database/sql"
	"github.com/VladislavSCV/api/middleware"
	"github.com/VladislavSCV/api/rest/handlers"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, db *sql.DB) {
	// Аутентификация
	auth := router.Group("/api/auth")
	{
		// Применение rate limiting к маршрутам
		auth.Use(middleware.RateLimiterMiddleware())
		auth.POST("/login", handlers.Login(db))
		auth.POST("/registration", handlers.Registration(db))
		auth.POST("/verify", handlers.Verify(db))
		auth.GET("/", middleware.AuthMiddleware(db), handlers.GetCurrentUser(db))
	}
}
