package main

import (
	"database/sql"
	"fmt"
	"github.com/VladislavSCV/api/rest/routes"
	_ "github.com/VladislavSCV/docs" // Импортируйте сгенерированную документацию
	"github.com/VladislavSCV/internal/redis"
	"github.com/VladislavSCV/pkg"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
)

func main() {
	// Загрузка переменных окружения
	pkg.LoadEnv()

	connStr := os.Getenv("POSTGRES")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Проверка соединения с базой данных
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")
	fmt.Println("Database connection:", db)

	// Подключение к Redis
	connStrR := os.Getenv("REDIS_CONN_STR")
	if connStr == "" {
		log.Fatalf("Redis connection string is not set in environment variables")
	}

	client, err := redis.ConnToRedis(connStrR)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer func() {
		if err := redis.CloseRedis(); err != nil {
			log.Printf("Failed to close Redis connection: %v", err)
		}
	}()

	// Инициализация Gin
	r := gin.Default()

	// Настройка маршрутов
	routes.SetupAuthRoutes(r, db, client)
	routes.SetupUserRoutes(r, db, client)
	routes.SetupGroupRoutes(r, db, client)
	routes.SetupScheduleRoutes(r, db, client)
	routes.SetupGradeRoutes(r, db, client)
	routes.SetupAttendanceRoutes(r, db, client)

	// Тестовый маршрут
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Маршрут для Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
