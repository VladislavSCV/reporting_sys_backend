package core

import (
	"database/sql"
	"fmt"
	"github.com/VladislavSCV/internal/models"
)

func GetAllGroups(db *sql.DB) ([]models.Group, error) {
	var groups []models.Group

	rows, err := db.Query(`
        SELECT g.id, g.name, COUNT(u.id) AS student_count, g.created_at, g.updated_at
        FROM groups g
        LEFT JOIN users u ON g.id = u.group_id
        GROUP BY g.id, g.name
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch groups: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.StudentCount, &group.CreatedAt, &group.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan group: %v", err)
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over groups: %v", err)
	}

	return groups, nil
}

func GetGroupByID(db *sql.DB, groupID int) (*models.GroupDetail, error) {
	var group models.GroupDetail

	// Получаем основную информацию о группе
	err := db.QueryRow(`
        SELECT id, name, created_at, updated_at
        FROM groups
        WHERE id = $1
    `, groupID).Scan(&group.ID, &group.Name, &group.CreatedAt, &group.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("group not found")
		}
		return nil, fmt.Errorf("failed to fetch group: %v", err)
	}

	// Получаем список студентов в группе
	rows, err := db.Query(`
        SELECT id, first_name, middle_name, last_name, login, created_at, updated_at
        FROM users
        WHERE group_id = $1
    `, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch students: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var student models.User
		if err := rows.Scan(&student.ID, &student.FirstName, &student.MiddleName, &student.LastName, &student.Login, &student.CreatedAt, &student.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan student: %v", err)
		}
		group.Students = append(group.Students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over students: %v", err)
	}

	// Получаем расписание для группы
	scheduleRows, err := db.Query(`
        SELECT s.id, sub.name AS subject_name, t.first_name || ' ' || t.last_name AS teacher_name, 
               s.day_of_week, s.start_time, s.end_time, s.location
        FROM schedules s
        JOIN subjects sub ON s.subject_id = sub.id
        JOIN users t ON s.teacher_id = t.id
        WHERE s.group_id = $1
    `, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedule: %v", err)
	}
	defer scheduleRows.Close()

	for scheduleRows.Next() {
		var schedule models.Schedule
		if err := scheduleRows.Scan(&schedule.ID, &schedule.SubjectName, &schedule.TeacherName, &schedule.DayOfWeek, &schedule.StartTime, &schedule.EndTime, &schedule.Location); err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %v", err)
		}
		group.Schedule = append(group.Schedule, schedule)
	}

	if err := scheduleRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over schedule: %v", err)
	}

	return &group, nil
}

func CreateGroup(db *sql.DB, group models.Group) (int, error) {
	var groupID int

	err := db.QueryRow(`
        INSERT INTO groups (name, created_at, updated_at)
        VALUES ($1, NOW(), NOW())
        RETURNING id
    `, group.Name).Scan(&groupID)
	if err != nil {
		return 0, fmt.Errorf("failed to create group: %v", err)
	}

	return groupID, nil
}

func UpdateGroup(db *sql.DB, group models.Group) error {
	_, err := db.Exec(`
        UPDATE groups
        SET name = $1, updated_at = NOW()
        WHERE id = $2
    `, group.Name, group.ID)
	if err != nil {
		return fmt.Errorf("failed to update group: %v", err)
	}

	return nil
}

func DeleteGroup(db *sql.DB, groupID int) error {
	_, err := db.Exec("DELETE FROM groups WHERE id = $1", groupID)
	if err != nil {
		return fmt.Errorf("failed to delete group: %v", err)
	}

	return nil
}
