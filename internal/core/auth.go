package core

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VladislavSCV/internal/models"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/argon2"
	"time"
)

func RegisterUser(db *sql.DB, user models.User) (int, error) {
	var userID int

	// Проверка уникальности логина
	var existingID int
	err := db.QueryRow("SELECT id FROM users WHERE login = $1", user.Login).Scan(&existingID)
	if err == nil {
		return 0, fmt.Errorf("login already exists")
	} else if err != sql.ErrNoRows {
		return 0, fmt.Errorf("database error: %v", err)
	}

	// Логирование для отладки
	fmt.Printf("Generated salt: %s\n", user.Salt)
	fmt.Printf("Generated hash: %s\n", user.Password)

	// Вставка нового пользователя
	err = db.QueryRow(
		`INSERT INTO users (first_name, middle_name, last_name, role_id, group_id, login, password, salt, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()) 
		 RETURNING id`,
		user.FirstName, user.MiddleName, user.LastName, user.RoleID, user.GroupID, user.Login, user.Password, user.Salt,
	).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to register user: %v", err)
	}

	return userID, nil
}

func AuthenticateUser(db *sql.DB, cache *redis.Client, login, password string) (*models.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%s", login)

	// Проверяем кеш
	cachedUser, err := cache.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var user models.User
		if err := json.Unmarshal(cachedUser, &user); err == nil {
			return &user, nil
		}
	}

	// Если данных нет в кеше, запрашиваем из базы данных
	var user models.User
	row := db.QueryRow("SELECT id, login, password, role_id FROM users WHERE login = $1", login)
	if err := row.Scan(&user.ID, &user.Login, &user.Password, &user.RoleID); err != nil {
		return nil, fmt.Errorf("failed to authenticate user: %v", err)
	}

	// Сохраняем данные в кеше
	userJSON, err := json.Marshal(user)
	if err == nil {
		cache.Set(ctx, cacheKey, userJSON, 24*time.Hour) // Кешируем на 24 часа
	}

	return &user, nil
}

func CheckPasswordHash(password, hashStr, saltStr string) (bool, error) {
	// Проверяем входные данные
	if password == "" || saltStr == "" || hashStr == "" {
		return false, errors.New("invalid input: password, salt or hash is empty")
	}

	// Декодируем соль и хеш из базы данных
	salt, err := base64.StdEncoding.DecodeString(saltStr)
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	hashBytes, err := base64.StdEncoding.DecodeString(hashStr)
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	// Логирование для отладки
	fmt.Printf("Salt (decoded): %x\n", salt)
	fmt.Printf("Hash (decoded): %x\n", hashBytes)

	// Хешируем введённый пароль с той же солью
	newHash := argon2.IDKey([]byte(password), salt, 3, 32*1024, 4, 32)

	// Логирование для отладки
	fmt.Printf("New hash (generated): %x\n", newHash)

	// Сравниваем хеши
	if subtle.ConstantTimeCompare(hashBytes, newHash) != 1 {
		return false, errors.New("password mismatch")
	}

	return true, nil
}

func GetCurrentUser(db *sql.DB, cache *redis.Client, userID int) (*models.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%d", userID)

	// Проверяем кеш
	cachedUser, err := cache.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var user models.User
		if err := json.Unmarshal(cachedUser, &user); err == nil {
			return &user, nil
		}
	}

	// Если данных нет в кеше, запрашиваем из базы данных
	var user models.User
	row := db.QueryRow("SELECT id, login, first_name, last_name, role_id, group_id, created_at, updated_at FROM users WHERE id = $1", userID)
	if err := row.Scan(&user.ID, &user.Login, &user.FirstName, &user.LastName, &user.RoleID, &user.GroupID, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// Сохраняем данные в кеше
	userJSON, err := json.Marshal(user)
	if err == nil {
		cache.Set(ctx, cacheKey, userJSON, 24*time.Hour) // Кешируем на 24 часа
	}

	return &user, nil
}
