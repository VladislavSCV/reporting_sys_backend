package handlers

import (
	"database/sql"
	"github.com/VladislavSCV/internal/core"
	"github.com/VladislavSCV/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// GetGroups godoc
// @Summary Получить список всех групп
// @Description Возвращает список всех групп
// @Tags Groups
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Group "Успешный ответ"  example([{"id": 1, "name": "Group A"}, {"id": 2, "name": "Group B"}])
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/group [get]
func GetGroups(db *sql.DB, cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Получен запрос на получение списка всех групп")

		groups, err := core.GetAllGroups(db, cache)
		if err != nil {
			log.Printf("Ошибка при получении списка групп: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		log.Printf("Успешно получен список групп")
		c.JSON(http.StatusOK, groups)
	}
}

// GetGroupByID godoc
// @Summary Получить информацию о группе по её ID
// @Description Возвращает информацию о конкретной группе
// @Tags Groups
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "ID группы"  example(1)
// @Success 200 {object} models.GroupDetail "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Группа не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/group/{id} [get]
func GetGroupByID(db *sql.DB, cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("id")
		log.Printf("Получен запрос на получение информации о группе с ID: %s", groupID)

		// Преобразуем groupID в int
		groupIDInt, err := strconv.Atoi(groupID)
		if err != nil {
			log.Printf("Ошибка преобразования groupID в int: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid group ID"})
			return
		}

		// Валидация groupID
		if groupIDInt <= 0 {
			log.Printf("Некорректный groupID: %d", groupIDInt)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "group ID must be positive"})
			return
		}

		group, err := core.GetGroupByID(db, cache, groupIDInt)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Printf("Группа с ID %d не найдена", groupIDInt)
				c.JSON(http.StatusNotFound, ErrorResponse{Error: "group not found"})
			} else {
				log.Printf("Ошибка при получении информации о группе: %v", err)
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			}
			return
		}

		log.Printf("Успешно получена информация о группе с ID: %d", groupIDInt)
		c.JSON(http.StatusOK, group)
	}
}

// CreateGroup godoc
// @Summary Создать группу
// @Description Создаёт новую группу
// @Tags Groups
// @Accept  json
// @Produce  json
// @Param   group  body  models.Group  true  "Данные о группе"  example({"name": "Group A"})
// @Success 200 {object} SuccessResponse "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/group [post]
func CreateGroup(db *sql.DB, cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var group models.Group
		if err := c.ShouldBindJSON(&group); err != nil {
			log.Printf("Ошибка привязки JSON: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		// Валидация данных
		if group.Name == "" {
			log.Printf("Некорректное имя группы: %s", group.Name)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "group name is required"})
			return
		}

		groupID, err := core.CreateGroup(db, cache, group)
		if err != nil {
			log.Printf("Ошибка при создании группы: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		log.Printf("Успешно создана группа с ID: %d", groupID)
		c.JSON(http.StatusOK, SuccessResponse{
			Message: "Group created successfully",
			Data:    gin.H{"group_id": groupID},
		})
	}
}

// UpdateGroup godoc
// @Summary Обновить информацию о группе
// @Description Обновляет информацию о существующей группе
// @Tags Groups
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "ID группы"  example(1)
// @Param   group  body  models.Group  true  "Обновлённые данные о группе"  example({"name": "Updated Group A"})
// @Success 200 {object} SuccessResponse "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/group/{id} [put]
func UpdateGroup(db *sql.DB, cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("id")
		log.Printf("Получен запрос на обновление группы с ID: %s", groupID)

		// Преобразуем groupID в int
		groupIDInt, err := strconv.Atoi(groupID)
		if err != nil {
			log.Printf("Ошибка преобразования groupID в int: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid group ID"})
			return
		}

		// Валидация groupID
		if groupIDInt <= 0 {
			log.Printf("Некорректный groupID: %d", groupIDInt)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "group ID must be positive"})
			return
		}

		var group models.Group
		if err := c.ShouldBindJSON(&group); err != nil {
			log.Printf("Ошибка привязки JSON: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		// Валидация данных
		if group.Name == "" {
			log.Printf("Некорректное имя группы: %s", group.Name)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "group name is required"})
			return
		}

		// Устанавливаем ID группы из параметра запроса
		group.ID = groupIDInt

		if err := core.UpdateGroup(db, cache, group); err != nil {
			log.Printf("Ошибка при обновлении группы: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		log.Printf("Успешно обновлена группа с ID: %d", groupIDInt)
		c.JSON(http.StatusOK, SuccessResponse{Message: "Group updated successfully"})
	}
}

// DeleteGroup godoc
// @Summary Удалить группу
// @Description Удаляет группу по её ID
// @Tags Groups
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "ID группы"  example(1)
// @Success 200 {object} SuccessResponse "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/group/{id} [delete]
func DeleteGroup(db *sql.DB, cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("id")
		log.Printf("Получен запрос на удаление группы с ID: %s", groupID)

		// Преобразуем groupID в int
		groupIDInt, err := strconv.Atoi(groupID)
		if err != nil {
			log.Printf("Ошибка преобразования groupID в int: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid group ID"})
			return
		}

		// Валидация groupID
		if groupIDInt <= 0 {
			log.Printf("Некорректный groupID: %d", groupIDInt)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "group ID must be positive"})
			return
		}

		if err := core.DeleteGroup(db, cache, groupIDInt); err != nil {
			log.Printf("Ошибка при удалении группы: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		log.Printf("Успешно удалена группа с ID: %d", groupIDInt)
		c.JSON(http.StatusOK, SuccessResponse{Message: "Group deleted successfully"})
	}
}
