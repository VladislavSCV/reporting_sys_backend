package core

import (
	"database/sql"
	"fmt"
	"github.com/VladislavSCV/internal/models"
)

func GetGradesByStudentID(db *sql.DB, studentID int) ([]models.GradeDetail, error) {
	var grades []models.GradeDetail

	rows, err := db.Query(`
        SELECT g.id, 
               g.student_id, 
               s.name AS subject_name, 
               g.value, 
               g.date, 
               g.created_at, 
               g.updated_at
        FROM grades g
        JOIN subjects s ON g.subject_id = s.id
        WHERE g.student_id = $1
    `, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch grades: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var grade models.GradeDetail
		if err := rows.Scan(
			&grade.ID,
			&grade.StudentID,
			&grade.SubjectName,
			&grade.Value,
			&grade.Date,
			&grade.CreatedAt,
			&grade.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan grade: %v", err)
		}
		grades = append(grades, grade)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over grades: %v", err)
	}

	return grades, nil
}

func GetGradesByGroupID(db *sql.DB, groupID int) ([]models.GradeDetail, error) {
	var grades []models.GradeDetail

	rows, err := db.Query(`
        SELECT g.id, 
               g.student_id, 
               u.first_name || ' ' || u.last_name AS student_name, 
               s.name AS subject_name, 
               g.value, 
               g.date, 
               g.created_at, 
               g.updated_at
        FROM grades g
        JOIN users u ON g.student_id = u.id
        JOIN subjects s ON g.subject_id = s.id
        WHERE u.group_id = $1
    `, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch grades: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var grade models.GradeDetail
		if err := rows.Scan(
			&grade.ID,
			&grade.StudentID,
			&grade.StudentName,
			&grade.SubjectName,
			&grade.Value,
			&grade.Date,
			&grade.CreatedAt,
			&grade.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan grade: %v", err)
		}
		grades = append(grades, grade)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over grades: %v", err)
	}

	return grades, nil
}

func CreateGrade(db *sql.DB, grade models.Grade) (int, error) {
	var gradeID int

	err := db.QueryRow(`
        INSERT INTO grades (student_id, subject_id, value, date, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id
    `, grade.StudentID, grade.SubjectID, grade.Value, grade.Date).Scan(&gradeID)
	if err != nil {
		return 0, fmt.Errorf("failed to create grade: %v", err)
	}

	return gradeID, nil
}

func UpdateGrade(db *sql.DB, grade models.Grade) error {
	_, err := db.Exec(`
        UPDATE grades
        SET student_id = $1, subject_id = $2, value = $3, date = $4, updated_at = NOW()
        WHERE id = $5
    `, grade.StudentID, grade.SubjectID, grade.Value, grade.Date, grade.ID)
	if err != nil {
		return fmt.Errorf("failed to update grade: %v", err)
	}

	return nil
}

func DeleteGrade(db *sql.DB, gradeID int) error {
	_, err := db.Exec("DELETE FROM grades WHERE id = $1", gradeID)
	if err != nil {
		return fmt.Errorf("failed to delete grade: %v", err)
	}

	return nil
}
