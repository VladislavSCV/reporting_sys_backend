package handlers

import (
	"database/sql"
	"github.com/VladislavSCV/internal/core"
	"github.com/VladislavSCV/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// GetSchedules godoc
// @Summary Получить общее расписание
// @Description Возвращает список всех занятий
// @Tags Schedules
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Schedule "Успешный ответ"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/schedule [get]
func GetSchedules(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Получен запрос на получение общего расписания")

		schedules, err := core.GetAllSchedules(db)
		if err != nil {
			log.Printf("Ошибка при получении расписания: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		log.Printf("Успешно получено общее расписание")
		c.JSON(http.StatusOK, schedules)
	}
}

// GetScheduleByID godoc
// @Summary Получить расписание по ID группы/преподавателя
// @Description Возвращает расписание для конкретной группы или преподавателя
// @Tags Schedules
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "ID группы или преподавателя"  example(1)
// @Success 200 {array} models.Schedule "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Расписание не найдено"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/schedule/{id} [get]
func GetScheduleByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("Получен запрос на получение расписания для ID: %s", id)

		// Преобразуем id в int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			log.Printf("Ошибка преобразования ID в int: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid ID"})
			return
		}

		// Валидация ID
		if idInt <= 0 {
			log.Printf("Некорректный ID: %d", idInt)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "ID must be positive"})
			return
		}

		schedules, err := core.GetScheduleByID(db, idInt)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Printf("Расписание для ID %d не найдено", idInt)
				c.JSON(http.StatusNotFound, ErrorResponse{Error: "schedule not found"})
			} else {
				log.Printf("Ошибка при получении расписания: %v", err)
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			}
			return
		}

		log.Printf("Успешно получено расписание для ID: %d", idInt)
		c.JSON(http.StatusOK, schedules)
	}
}

// CreateSchedule godoc
// @Summary Создать занятие
// @Description Создаёт новое занятие в расписании
// @Tags Schedules
// @Accept  json
// @Produce  json
// @Param   schedule  body  models.Schedule  true  "Данные о занятии"  example({"group_id": 1, "subject_id": 1, "teacher_id": 1, "day_of_week": 1, "start_time": "09:00", "end_time": "10:30", "location": "Room 101"})
// @Success 200 {object} SuccessResponse "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/schedule [post]
func CreateSchedule(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var schedule models.Schedule
		if err := c.ShouldBindJSON(&schedule); err != nil {
			log.Printf("Ошибка привязки JSON: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		// Валидация данных
		if schedule.GroupID <= 0 {
			log.Printf("Некорректный group_id: %d", schedule.GroupID)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "group_id must be positive"})
			return
		}
		if schedule.SubjectID <= 0 {
			log.Printf("Некорректный subject_id: %d", schedule.SubjectID)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "subject_id must be positive"})
			return
		}
		if schedule.TeacherID <= 0 {
			log.Printf("Некорректный teacher_id: %d", schedule.TeacherID)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "teacher_id must be positive"})
			return
		}
		if schedule.DayOfWeek < 1 || schedule.DayOfWeek > 7 {
			log.Printf("Некорректный день недели: %d", schedule.DayOfWeek)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "day_of_week must be between 1 and 7"})
			return
		}
		if schedule.StartTime == "" || schedule.EndTime == "" {
			log.Printf("Некорректное время начала или окончания занятия")
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "start_time and end_time must be provided"})
			return
		}
		if schedule.Location == "" {
			log.Printf("Некорректное местоположение")
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "location must be provided"})
			return
		}

		scheduleID, err := core.CreateSchedule(db, schedule)
		if err != nil {
			log.Printf("Ошибка при создании занятия: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		log.Printf("Успешно создано занятие с ID: %d", scheduleID)
		c.JSON(http.StatusOK, SuccessResponse{
			Message: "Schedule created successfully",
			Data:    gin.H{"schedule_id": scheduleID},
		})
	}
}

// UpdateSchedule godoc
// @Summary Обновить занятие
// @Description Обновляет информацию о существующем занятии
// @Tags Schedules
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "ID занятия"  example(1)
// @Param   updates  body  map[string]interface{}  true  "Обновлённые данные о занятии"  example({"day_of_week": 2, "start_time": "10:00", "end_time": "11:30", "location": "Room 202"})
// @Success 200 {object} SuccessResponse "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/schedule/{id} [put]
func UpdateSchedule(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("Получен запрос на обновление занятия с ID: %s", id)

		// Преобразуем id в int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			log.Printf("Ошибка преобразования ID в int: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid ID"})
			return
		}

		// Валидация ID
		if idInt <= 0 {
			log.Printf("Некорректный ID: %d", idInt)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "ID must be positive"})
			return
		}

		var updates map[string]interface{}
		if err := c.ShouldBindJSON(&updates); err != nil {
			log.Printf("Ошибка привязки JSON: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		// Валидация дня недели, если он передан
		if dayOfWeek, ok := updates["day_of_week"].(float64); ok {
			if dayOfWeek < 1 || dayOfWeek > 7 {
				log.Printf("Некорректный день недели: %v", dayOfWeek)
				c.JSON(http.StatusBadRequest, ErrorResponse{Error: "day_of_week must be between 1 and 7"})
				return
			}
		}

		// Устанавливаем ID занятия из параметра запроса
		updates["id"] = idInt

		if err := core.UpdateSchedule(db, updates); err != nil {
			log.Printf("Ошибка при обновлении занятия: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		log.Printf("Успешно обновлено занятие с ID: %d", idInt)
		c.JSON(http.StatusOK, SuccessResponse{Message: "Schedule updated successfully"})
	}
}

// DeleteSchedule godoc
// @Summary Удалить занятие
// @Description Удаляет занятие по его ID
// @Tags Schedules
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "ID занятия"  example(1)
// @Success 200 {object} SuccessResponse "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/schedule/{id} [delete]
func DeleteSchedule(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("Получен запрос на удаление занятия с ID: %s", id)

		// Преобразуем id в int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			log.Printf("Ошибка преобразования ID в int: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid ID"})
			return
		}

		// Валидация ID
		if idInt <= 0 {
			log.Printf("Некорректный ID: %d", idInt)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "ID must be positive"})
			return
		}

		if err := core.DeleteSchedule(db, idInt); err != nil {
			log.Printf("Ошибка при удалении занятия: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		log.Printf("Успешно удалено занятие с ID: %d", idInt)
		c.JSON(http.StatusOK, SuccessResponse{Message: "Schedule deleted successfully"})
	}
}
