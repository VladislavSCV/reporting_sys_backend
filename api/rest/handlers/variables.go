package handlers

import (
	"github.com/VladislavSCV/internal/database"
	"time"
)

var db = database.DB

// SuccessResponse представляет успешный ответ API.
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // Опциональные данные
}

// ErrorResponse представляет ошибочный ответ API.
type ErrorResponse struct {
	Error string `json:"error"`
}

// LoginResponse представляет ответ на успешный вход в систему.
type LoginResponse struct {
	Token string `json:"token"`
}

// VerifyResponse представляет ответ на проверку токена.
type VerifyResponse struct {
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
}

// UserResponse представляет информацию о пользователе.
type UserResponse struct {
	ID         int       `json:"id"`
	FirstName  string    `json:"first_name"`
	MiddleName string    `json:"middle_name"`
	LastName   string    `json:"last_name"`
	RoleID     int       `json:"role_id"`
	GroupID    *int      `json:"group_id"`
	Login      string    `json:"login"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// GradeResponse представляет информацию об оценке.
type GradeResponse struct {
	ID        int       `json:"id"`
	StudentID int       `json:"student_id"`
	SubjectID int       `json:"subject_id"`
	Grade     int       `json:"grade"`
	Date      time.Time `json:"date"`
}
