package models

import "time"

type Grade struct {
	ID        int       `json:"id"`
	StudentID int       `json:"student_id"`
	SubjectID int       `json:"subject_id"`
	Value     int       `json:"value"` // Оценка (2-5)
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GradeDetail struct {
	ID          int       `json:"id"`
	StudentID   int       `json:"student_id"`
	StudentName string    `json:"student_name"`
	SubjectName string    `json:"subject_name"`
	Value       int       `json:"value"`
	Date        string    `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
