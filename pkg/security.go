package pkg

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)

type HashResult struct {
	Salt string
	Hash string
}

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var p params

type PasswordHash struct {
	Hash string
}

// GenerateSalt генерирует случайную соль
func GenerateSalt() (string, error) {
	salt := make([]byte, 16) // длина соли
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// GenerateHashFromPassword хеширует пароль с солью и возвращает результат
func GenerateHashFromPassword(password, salt string) (*PasswordHash, error) {
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return nil, err
	}

	// Хешируем пароль с солью
	hash := argon2.IDKey([]byte(password), saltBytes, 3, 32*1024, 4, 32)
	return &PasswordHash{
		Hash: base64.StdEncoding.EncodeToString(hash),
	}, nil
}

// VerifyPassword проверяет, соответствует ли введенный пароль хешу
func VerifyPassword(password, saltStr, hashStr string) (bool, error) {
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

	// Хешируем введённый пароль с той же солью
	newHash := argon2.IDKey([]byte(password), salt, 3, 32*1024, 4, 32)

	// Сравниваем хеши
	if subtle.ConstantTimeCompare(hashBytes, newHash) != 1 {
		return false, errors.New("password mismatch")
	}

	return true, nil
}

// HashPasswordWithSalt хеширует пароль с переданной солью
func HashPasswordWithSalt(password string, salt string) (string, error) {
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}

	// Используем argon2 для хеширования пароля с солью
	hash := argon2.IDKey([]byte(password), saltBytes, 3, 32*1024, 4, 32)

	return base64.StdEncoding.EncodeToString(hash), nil
}

// CreateHashWithSalt генерирует хеш пароля и возвращает соль вместе с хешем
func CreateHashWithSalt(password string) (*HashResult, error) {
	// Генерация соли
	salt, err := GenerateSalt()
	if err != nil {
		return nil, err
	}

	// Генерация хеша с солью С САХАРОМ!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	hashResult, err := GenerateHashFromPassword(password, salt)
	if err != nil {
		return nil, err
	}

	// Возвращаем соль и хеш
	return &HashResult{
		Salt: salt,
		Hash: hashResult.Hash,
	}, nil
}
