package middleware

import (
	"github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"time"
)

func RateLimiterMiddleware() gin.HandlerFunc {
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Minute, // Ограничение: 10 запросов в минуту
		Limit: 10,
	})

	return ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: func(c *gin.Context, info ratelimit.Info) {
			c.JSON(429, gin.H{
				"error":   "Too many requests",
				"message": "Try again in " + time.Until(info.ResetTime).String(),
			})
		},
		KeyFunc: func(c *gin.Context) string {
			ip := c.ClientIP()
			println("Client IP:", ip) // Логируем IP-адрес
			return ip
		},
	})
}
