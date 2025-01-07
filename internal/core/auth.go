package core

import (
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/VladislavSCV/internal/models"
	"golang.org/x/crypto/argon2"
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

func AuthenticateUser(db *sql.DB, login, password string) (*models.User, error) {
	var user models.User

	// Получение данных пользователя
	err := db.QueryRow(
		`SELECT id, first_name, middle_name, last_name, role_id, group_id, login, password, salt, created_at, updated_at 
		 FROM users 
		 WHERE login = $1`,
		login,
	).Scan(
		&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.RoleID, &user.GroupID,
		&user.Login, &user.Password, &user.Salt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}

	// Проверка пароля
	if isValid, err := CheckPasswordHash(password, user.Password, user.Salt); err != nil || !isValid {
		return nil, fmt.Errorf("invalid password")
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

func GetCurrentUser(db *sql.DB, userID int) (*models.User, error) {
	var user models.User

	err := db.QueryRow(
		"SELECT id, first_name, middle_name, last_name, role_id, group_id, login, created_at, updated_at FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.RoleID, &user.GroupID, &user.Login, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
