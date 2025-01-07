package core

import (
	"database/sql"
	"fmt"
	"github.com/VladislavSCV/internal/models"
)

func GetAttendanceByStudentID(db *sql.DB, studentID int) ([]models.AttendanceDetail, error) {
	var attendances []models.AttendanceDetail

	rows, err := db.Query(`
        SELECT a.id, 
               a.student_id, 
               s.name AS subject_name, 
               a.date, 
               a.status, 
               a.created_at, 
               a.updated_at
        FROM attendance a
        JOIN subjects s ON a.subject_id = s.id
        WHERE a.student_id = $1
    `, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attendance: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var attendance models.AttendanceDetail
		if err := rows.Scan(
			&attendance.ID,
			&attendance.StudentID,
			&attendance.SubjectName,
			&attendance.Date,
			&attendance.Status,
			&attendance.CreatedAt,
			&attendance.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan attendance: %v", err)
		}
		attendances = append(attendances, attendance)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over attendance: %v", err)
	}

	return attendances, nil
}

func GetAttendanceByGroupID(db *sql.DB, groupID int) ([]models.AttendanceDetail, error) {
	var attendances []models.AttendanceDetail

	rows, err := db.Query(`
        SELECT a.id, 
               a.student_id, 
               u.first_name || ' ' || u.last_name AS student_name, 
               s.name AS subject_name, 
               a.date, 
               a.status, 
               a.created_at, 
               a.updated_at
        FROM attendance a
        JOIN users u ON a.student_id = u.id
        JOIN subjects s ON a.subject_id = s.id
        WHERE u.group_id = $1
    `, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attendance: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var attendance models.AttendanceDetail
		if err := rows.Scan(
			&attendance.ID,
			&attendance.StudentID,
			&attendance.StudentName,
			&attendance.SubjectName,
			&attendance.Date,
			&attendance.Status,
			&attendance.CreatedAt,
			&attendance.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan attendance: %v", err)
		}
		attendances = append(attendances, attendance)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over attendance: %v", err)
	}

	return attendances, nil
}

func CreateAttendance(db *sql.DB, attendance models.Attendance) (int, error) {
	var attendanceID int

	err := db.QueryRow(`
        INSERT INTO attendance (student_id, subject_id, date, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id
    `, attendance.StudentID, attendance.SubjectID, attendance.Date, attendance.Status).Scan(&attendanceID)
	if err != nil {
		return 0, fmt.Errorf("failed to create attendance: %v", err)
	}

	return attendanceID, nil
}

func UpdateAttendance(db *sql.DB, attendance models.Attendance) error {
	_, err := db.Exec(`
        UPDATE attendance
        SET student_id = $1, subject_id = $2, date = $3, status = $4, updated_at = NOW()
        WHERE id = $5
    `, attendance.StudentID, attendance.SubjectID, attendance.Date, attendance.Status, attendance.ID)
	if err != nil {
		return fmt.Errorf("failed to update attendance: %v", err)
	}

	return nil
}
