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

// GetGradesByStudentID godoc
// @Summary Получить оценки студента по его ID
// @Description Возвращает список оценок для конкретного студента
// @Tags Grades
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "ID студента"  example(1)
// @Success 200 {array} models.GradeDetail "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Оценки не найдены"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/grades/student/{id} [get]
func GetGradesByStudentID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := c.Param("id")
		log.Printf("Получен запрос на получение оценок для студента с ID: %s", studentID)

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

		grades, err := core.GetGradesByStudentID(db, studentIDInt)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Printf("Оценки для студента с ID %d не найдены", studentIDInt)
				c.JSON(http.StatusNotFound, ErrorResponse{Error: "grades not found"})
			} else {
				log.Printf("Ошибка при получении оценок: %v", err)
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			}
			return
		}

		log.Printf("Успешно получены оценки для студента с ID: %d", studentIDInt)
		c.JSON(http.StatusOK, grades)
	}
}

// GetGradesByGroupID godoc
// @Summary Получить оценки группы по её ID
// @Description Возвращает список оценок для конкретной группы
// @Tags Grades
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "ID группы"  example(1)
// @Success 200 {array} models.GradeDetail "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Оценки не найдены"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/grades/group/{id} [get]
func GetGradesByGroupID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("id")
		log.Printf("Получен запрос на получение оценок для группы с ID: %s", groupID)

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

		grades, err := core.GetGradesByGroupID(db, groupIDInt)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Printf("Оценки для группы с ID %d не найдены", groupIDInt)
				c.JSON(http.StatusNotFound, ErrorResponse{Error: "grades not found"})
			} else {
				log.Printf("Ошибка при получении оценок: %v", err)
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			}
			return
		}

		log.Printf("Успешно получены оценки для группы с ID: %d", groupIDInt)
		c.JSON(http.StatusOK, grades)
	}
}

// CreateGrade godoc
// @Summary Создать оценку
// @Description Создаёт новую запись об оценке
// @Tags Grades
// @Accept  json
// @Produce  json
// @Param   grade  body  models.Grade  true  "Данные об оценке"  example({"student_id": 1, "subject_id": 1, "value": 5, "date": "2023-10-01T00:00:00Z"})
// @Success 200 {object} SuccessResponse "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/grades [post]
func CreateGrade(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var grade models.Grade
		if err := c.ShouldBindJSON(&grade); err != nil {
			log.Printf("Ошибка привязки JSON: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		// Валидация данных
		if grade.StudentID <= 0 {
			log.Printf("Некорректный student_id: %d", grade.StudentID)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "student_id must be positive"})
			return
		}
		if grade.SubjectID <= 0 {
			log.Printf("Некорректный subject_id: %d", grade.SubjectID)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "subject_id must be positive"})
			return
		}
		if grade.Value < 2 || grade.Value > 5 {
			log.Printf("Некорректная оценка: %d", grade.Value)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "value must be between 2 and 5"})
			return
		}
		if grade.Date.IsZero() {
			log.Printf("Некорректная дата: %v", grade.Date)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "date must be provided"})
			return
		}

		gradeID, err := core.CreateGrade(db, grade)
		if err != nil {
			log.Printf("Ошибка при создании оценки: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		log.Printf("Успешно создана оценка с ID: %d", gradeID)
		c.JSON(http.StatusOK, SuccessResponse{
			Message: "Grade created successfully",
			Data:    gin.H{"grade_id": gradeID},
		})
	}
}

// UpdateGrade godoc
// @Summary Обновить оценку
// @Description Обновляет существующую запись об оценке
// @Tags Grades
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "ID оценки"  example(1)
// @Param   grade  body  models.Grade  true  "Обновлённые данные об оценке"  example({"student_id": 1, "subject_id": 1, "value": 4, "date": "2023-10-01T00:00:00Z"})
// @Success 200 {object} SuccessResponse "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/grades/{id} [put]
func UpdateGrade(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("Получен запрос на обновление оценки с ID: %s", id)

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

		var grade models.Grade
		if err := c.ShouldBindJSON(&grade); err != nil {
			log.Printf("Ошибка привязки JSON: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		// Валидация данных
		if grade.StudentID <= 0 {
			log.Printf("Некорректный student_id: %d", grade.StudentID)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "student_id must be positive"})
			return
		}
		if grade.SubjectID <= 0 {
			log.Printf("Некорректный subject_id: %d", grade.SubjectID)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "subject_id must be positive"})
			return
		}
		if grade.Value < 2 || grade.Value > 5 {
			log.Printf("Некорректная оценка: %d", grade.Value)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "value must be between 2 and 5"})
			return
		}
		if grade.Date.IsZero() {
			log.Printf("Некорректная дата: %v", grade.Date)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "date must be provided"})
			return
		}

		// Устанавливаем ID оценки из параметра запроса
		grade.ID = idInt

		if err := core.UpdateGrade(db, grade); err != nil {
			log.Printf("Ошибка при обновлении оценки: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		log.Printf("Успешно обновлена оценка с ID: %d", idInt)
		c.JSON(http.StatusOK, SuccessResponse{Message: "Grade updated successfully"})
	}
}

// DeleteGrade godoc
// @Summary Удалить оценку
// @Description Удаляет запись об оценке
// @Tags Grades
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "ID оценки"  example(1)
// @Success 200 {object} SuccessResponse "Успешный ответ"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/grades/{id} [delete]
func DeleteGrade(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("Получен запрос на удаление оценки с ID: %s", id)

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

		if err := core.DeleteGrade(db, idInt); err != nil {
			log.Printf("Ошибка при удалении оценки: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		log.Printf("Успешно удалена оценка с ID: %d", idInt)
		c.JSON(http.StatusOK, SuccessResponse{Message: "Grade deleted successfully"})
	}
}
