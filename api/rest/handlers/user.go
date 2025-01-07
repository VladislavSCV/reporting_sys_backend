package handlers

import (
	"database/sql"
	"github.com/VladislavSCV/internal/core"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// GetUsers godoc
// @Summary Получить список всех пользователей
// @Description Возвращает список всех пользователей
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User "Успешный ответ"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/user [get]
func GetUsers(db *sql.DB, cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Получен запрос на получение списка всех пользователей")

		users, err := core.GetAllUsers(db, cache)
		if err != nil {
			log.Printf("Ошибка при получении списка пользователей: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		log.Printf("Успешно получен список всех пользователей")
		c.JSON(http.StatusOK, users)
	}
}

// GetStudents godoc
// @Summary Получить список студентов
// @Description Возвращает список всех студентов
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User "Успешный ответ"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/user/students [get]
func GetStudents(db *sql.DB, cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Получен запрос на получение списка студентов")

		students, err := core.GetStudents(db, cache)
		if err != nil {
			log.Printf("Ошибка при получении списка студентов: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		log.Printf("Успешно получен список студентов")
		c.JSON(http.StatusOK, students)
	}
}

// GetTeachers godoc
// @Summary Получить список преподавателей
// @Description Возвращает список всех преподавателей
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User "Успешный ответ"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/user/teachers [get]
func GetTeachers(db *sql.DB, cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Получен запрос на получение списка преподавателей")

		teachers, err := core.GetTeachers(db, cache)
		if err != nil {
			log.Printf("Ошибка при получении списка преподавателей: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		log.Printf("Успешно получен список преподавателей")
		c.JSON(http.StatusOK, teachers)
	}
}

// GetUserByID godoc
// @Summary Получить информацию о пользователе по его ID
// @Description Возвращает информацию о конкретном пользователе
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "ID пользователя"  example(1)
// @Success 200 {object} models.User "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/user/{id} [get]
func GetUserByID(db *sql.DB, cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		log.Printf("Получен запрос на получение информации о пользователе с ID: %s", userID)

		// Преобразуем userID в int
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			log.Printf("Ошибка преобразования userID в int: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user ID"})
			return
		}

		// Валидация userID
		if userIDInt <= 0 {
			log.Printf("Некорректный userID: %d", userIDInt)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "user ID must be positive"})
			return
		}

		user, err := core.GetUserByID(db, cache, userIDInt)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Printf("Пользователь с ID %d не найден", userIDInt)
				c.JSON(http.StatusNotFound, ErrorResponse{Error: "user not found"})
			} else {
				log.Printf("Ошибка при получении информации о пользователе: %v", err)
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			}
			return
		}

		log.Printf("Успешно получена информация о пользователе с ID: %d", userIDInt)
		c.JSON(http.StatusOK, user)
	}
}

// UpdateUser godoc
// @Summary Обновить информацию о пользователе
// @Description Обновляет информацию о существующем пользователе
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "ID пользователя"  example(1)
// @Param   updates  body  map[string]interface{}  true  "Обновлённые данные о пользователе"  example({"first_name": "John", "last_name": "Doe"})
// @Success 200 {object} SuccessResponse "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/user/{id} [put]
func UpdateUser(db *sql.DB, cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		log.Printf("Получен запрос на обновление пользователя с ID: %s", userID)

		// Преобразуем userID в int
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			log.Printf("Ошибка преобразования userID в int: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user ID"})
			return
		}

		// Валидация userID
		if userIDInt <= 0 {
			log.Printf("Некорректный userID: %d", userIDInt)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "user ID must be positive"})
			return
		}

		var updates map[string]interface{}
		if err := c.ShouldBindJSON(&updates); err != nil {
			log.Printf("Ошибка привязки JSON: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		// Валидация обновлений
		if firstName, ok := updates["first_name"].(string); ok && firstName == "" {
			log.Printf("Некорректное имя пользователя: %s", firstName)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "first_name cannot be empty"})
			return
		}
		if lastName, ok := updates["last_name"].(string); ok && lastName == "" {
			log.Printf("Некорректная фамилия пользователя: %s", lastName)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "last_name cannot be empty"})
			return
		}

		// Устанавливаем ID пользователя из параметра запроса
		updates["id"] = userIDInt

		if err := core.UpdateUser(db, cache, updates); err != nil {
			log.Printf("Ошибка при обновлении пользователя: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		log.Printf("Успешно обновлен пользователь с ID: %d", userIDInt)
		c.JSON(http.StatusOK, SuccessResponse{Message: "User updated successfully"})
	}
}

// DeleteUser godoc
// @Summary Удалить пользователя
// @Description Удаляет пользователя по его ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "ID пользователя"  example(1)
// @Success 200 {object} SuccessResponse "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/user/{id} [delete]
func DeleteUser(db *sql.DB, cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		log.Printf("Получен запрос на удаление пользователя с ID: %s", userID)

		// Преобразуем userID в int
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			log.Printf("Ошибка преобразования userID в int: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user ID"})
			return
		}

		// Валидация userID
		if userIDInt <= 0 {
			log.Printf("Некорректный userID: %d", userIDInt)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "user ID must be positive"})
			return
		}

		if err := core.DeleteUser(db, cache, userIDInt); err != nil {
			log.Printf("Ошибка при удалении пользователя: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		log.Printf("Успешно удален пользователь с ID: %d", userIDInt)
		c.JSON(http.StatusOK, SuccessResponse{Message: "User deleted successfully"})
	}
}
