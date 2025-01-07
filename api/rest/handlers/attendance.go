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

// GetAttendanceByStudentID godoc
// @Summary Получить посещаемость студента по его ID
// @Description Возвращает список посещаемости для конкретного студента
// @Tags Attendance
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "ID студента"  example(1)
// @Success 200 {array} models.Attendance "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Посещаемость не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/attendance/student/{id} [get]
func GetAttendanceByStudentID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := c.Param("id")
		log.Printf("Получен запрос на получение посещаемости для студента с ID: %s", studentID)

		// Преобразуем studentID в int
		studentIDInt, err := strconv.Atoi(studentID)
		if err != nil {
			log.Printf("Ошибка преобразования studentID в int: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid student ID"})
			return
		}

		// Валидация studentID
		if studentIDInt <= 0 {
			log.Printf("Некорректный studentID: %d", studentIDInt)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "student ID must be positive"})
			return
		}

		attendances, err := core.GetAttendanceByStudentID(db, studentIDInt)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Printf("Посещаемость для студента с ID %d не найдена", studentIDInt)
				c.JSON(http.StatusNotFound, ErrorResponse{Error: "attendance not found"})
			} else {
				log.Printf("Ошибка при получении посещаемости: %v", err)
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			}
			return
		}

		log.Printf("Успешно получена посещаемость для студента с ID: %d", studentIDInt)
		c.JSON(http.StatusOK, attendances)
	}
}

// GetAttendanceByGroupID godoc
// @Summary Получить посещаемость группы по её ID
// @Description Возвращает список посещаемости для конкретной группы
// @Tags Attendance
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "ID группы"  example(1)
// @Success 200 {array} models.Attendance "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Посещаемость не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/attendance/group/{id} [get]
func GetAttendanceByGroupID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("id")
		log.Printf("Получен запрос на получение посещаемости для группы с ID: %s", groupID)

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

		attendances, err := core.GetAttendanceByGroupID(db, groupIDInt)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Printf("Посещаемость для группы с ID %d не найдена", groupIDInt)
				c.JSON(http.StatusNotFound, ErrorResponse{Error: "attendance not found"})
			} else {
				log.Printf("Ошибка при получении посещаемости: %v", err)
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			}
			return
		}

		log.Printf("Успешно получена посещаемость для группы с ID: %d", groupIDInt)
		c.JSON(http.StatusOK, attendances)
	}
}

// CreateAttendance godoc
// @Summary Создать отметку посещаемости
// @Description Создаёт новую запись о посещаемости
// @Tags Attendance
// @Accept  json
// @Produce  json
// @Param   attendance  body  models.Attendance  true  "Данные о посещаемости"  example({"student_id": 1, "subject_id": 1, "date": "2023-10-01T00:00:00Z", "status": "present"})
// @Success 200 {object} SuccessResponse "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/attendance [post]
func CreateAttendance(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var attendance models.Attendance
		if err := c.ShouldBindJSON(&attendance); err != nil {
			log.Printf("Ошибка привязки JSON: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		// Валидация данных
		if attendance.StudentID <= 0 {
			log.Printf("Некорректный student_id: %d", attendance.StudentID)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "student_id must be positive"})
			return
		}
		if attendance.SubjectID <= 0 {
			log.Printf("Некорректный subject_id: %d", attendance.SubjectID)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "subject_id must be positive"})
			return
		}
		if attendance.Date.IsZero() {
			log.Printf("Некорректная дата: %v", attendance.Date)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "date must be provided"})
			return
		}
		if attendance.Status != "present" && attendance.Status != "absent" {
			log.Printf("Некорректный статус: %s", attendance.Status)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "status must be 'present' or 'absent'"})
			return
		}

		attendanceID, err := core.CreateAttendance(db, attendance)
		if err != nil {
			log.Printf("Ошибка при создании посещаемости: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		log.Printf("Успешно создана посещаемость с ID: %d", attendanceID)
		c.JSON(http.StatusOK, SuccessResponse{
			Message: "Attendance created successfully",
			Data:    gin.H{"attendance_id": attendanceID},
		})
	}
}

// UpdateAttendance godoc
// @Summary Обновить отметку посещаемости
// @Description Обновляет существующую запись о посещаемости
// @Tags Attendance
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "ID посещаемости"  example(1)
// @Param   attendance  body  models.Attendance  true  "Обновлённые данные о посещаемости"  example({"student_id": 1, "subject_id": 1, "date": "2023-10-01T00:00:00Z", "status": "absent"})
// @Success 200 {object} SuccessResponse "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/attendance/{id} [put]
func UpdateAttendance(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("Получен запрос на обновление посещаемости с ID: %s", id)

		// Преобразуем id в int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			log.Printf("Ошибка преобразования id в int: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid ID"})
			return
		}

		// Валидация id
		if idInt <= 0 {
			log.Printf("Некорректный id: %d", idInt)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "ID must be positive"})
			return
		}

		var attendance models.Attendance
		if err := c.ShouldBindJSON(&attendance); err != nil {
			log.Printf("Ошибка привязки JSON: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		// Валидация данных
		if attendance.StudentID <= 0 {
			log.Printf("Некорректный student_id: %d", attendance.StudentID)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "student_id must be positive"})
			return
		}
		if attendance.SubjectID <= 0 {
			log.Printf("Некорректный subject_id: %d", attendance.SubjectID)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "subject_id must be positive"})
			return
		}
		if attendance.Date.IsZero() {
			log.Printf("Некорректная дата: %v", attendance.Date)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "date must be provided"})
			return
		}
		if attendance.Status != "present" && attendance.Status != "absent" {
			log.Printf("Некорректный статус: %s", attendance.Status)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "status must be 'present' or 'absent'"})
			return
		}

		// Устанавливаем ID отметки посещаемости из параметра запроса
		attendance.ID = idInt

		if err := core.UpdateAttendance(db, attendance); err != nil {
			log.Printf("Ошибка при обновлении посещаемости: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		log.Printf("Успешно обновлена посещаемость с ID: %d", idInt)
		c.JSON(http.StatusOK, SuccessResponse{Message: "Attendance updated successfully"})
	}
}
