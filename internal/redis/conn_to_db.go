package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var (
	ctx context.Context
	rdb *redis.Client
)

// ConnToRedis устанавливает соединение с Redis.
// Принимает строку подключения (например, "rediss://:<password>@<host>:<port>/<db>").
// Возвращает клиент Redis и ошибку, если подключение не удалось.
func ConnToRedis(connStr string) (*redis.Client, error) {
	// Парсим строку подключения
	opt, err := redis.ParseURL(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis connection string: %v", err)
	}

	// Создаем клиент Redis
	rdb = redis.NewClient(opt)

	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	log.Println("Successfully connected to Redis")
	return rdb, nil
}

// GetClient возвращает клиент Redis для использования в других пакетах.
// Если клиент не инициализирован, возвращает nil.
func GetClient() *redis.Client {
	return rdb
}

// CloseRedis закрывает соединение с Redis.
func CloseRedis() error {
	if rdb != nil {
		err := rdb.Close()
		if err != nil {
			return fmt.Errorf("failed to close Redis connection: %v", err)
		}
		log.Println("Redis connection closed")
		rdb = nil // Обнуляем клиент
	}
	return nil
}
