package core

import (
	"database/sql"
	"fmt"
	"github.com/VladislavSCV/internal/models"
)

func GetAllSchedules(db *sql.DB) ([]models.Schedule, error) {
	var schedules []models.Schedule

	rows, err := db.Query(`
        SELECT s.id, 
               g.name AS group_name, 
               sub.name AS subject_name, 
               t.first_name || ' ' || t.last_name AS teacher_name, 
               s.day_of_week, 
               s.start_time, 
               s.end_time, 
               s.location, 
               s.created_at, 
               s.updated_at,
			   sub.id AS subject_id,
			   g.id AS group_id,
               t.id AS teacher_id
        FROM schedules s
        JOIN groups g ON s.group_id = g.id
        JOIN subjects sub ON s.subject_id = sub.id
        JOIN users t ON s.teacher_id = t.id
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedules: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var schedule models.Schedule
		if err := rows.Scan(
			&schedule.ID,
			&schedule.GroupName,
			&schedule.SubjectName,
			&schedule.TeacherName,
			&schedule.DayOfWeek,
			&schedule.StartTime,
			&schedule.EndTime,
			&schedule.Location,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
			&schedule.SubjectID,
			&schedule.GroupID,
			&schedule.TeacherID,
		); err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %v", err)
		}
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over schedules: %v", err)
	}

	return schedules, nil
}

func GetScheduleByID(db *sql.DB, id int) ([]models.Schedule, error) {
	var schedules []models.Schedule

	rows, err := db.Query(`
        SELECT s.id, 
               g.name AS group_name, 
               sub.name AS subject_name, 
               t.first_name || ' ' || t.last_name AS teacher_name, 
               s.day_of_week, 
               s.start_time, 
               s.end_time, 
               s.location, 
               s.created_at, 
               s.updated_at,
               sub.id AS subject_id,
			   g.id AS group_id,
               t.id AS teacher_id
        FROM schedules s
        JOIN groups g ON s.group_id = g.id
        JOIN subjects sub ON s.subject_id = sub.id
        JOIN users t ON s.teacher_id = t.id
        WHERE s.group_id = $1
    `, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedule: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var schedule models.Schedule
		if err := rows.Scan(
			&schedule.ID,
			&schedule.GroupName,
			&schedule.SubjectName,
			&schedule.TeacherName,
			&schedule.DayOfWeek,
			&schedule.StartTime,
			&schedule.EndTime,
			&schedule.Location,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
			&schedule.SubjectID,
			&schedule.GroupID,
			&schedule.TeacherID,
		); err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %v", err)
		}
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over schedules: %v", err)
	}

	return schedules, nil
}

func CreateSchedule(db *sql.DB, schedule models.Schedule) (int, error) {
	var scheduleID int

	err := db.QueryRow(`
        INSERT INTO schedules (group_id, subject_id, teacher_id, day_of_week, start_time, end_time, location, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
        RETURNING id
    `, schedule.GroupID, schedule.SubjectID, schedule.TeacherID, schedule.DayOfWeek, schedule.StartTime, schedule.EndTime, schedule.Location).Scan(&scheduleID)
	if err != nil {
		return 0, fmt.Errorf("failed to create schedule: %v", err)
	}

	return scheduleID, nil
}

func UpdateSchedule(db *sql.DB, updates map[string]interface{}) error {
	query := "UPDATE schedules SET "
	var args []interface{}
	i := 1

	for k, v := range updates {
		query += fmt.Sprintf("%s = $%d, ", k, i) // Используем $ для подстановки
		args = append(args, v)
		i++
	}

	// Удаляем последнюю запятую
	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = $%d", i)
	args = append(args, updates["id"])

	// Выполнение запроса
	_, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

func DeleteSchedule(db *sql.DB, scheduleID int) error {
	_, err := db.Exec("DELETE FROM schedules WHERE id = $1", scheduleID)
	if err != nil {
		return fmt.Errorf("failed to delete schedule: %v", err)
	}

	return nil
}
